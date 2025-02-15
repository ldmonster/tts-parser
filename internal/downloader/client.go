package downloader

import (
	"bytes"
	"cmp"
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"sync"

	service "github.com/ldmonster/tts-parser/internal"
	"github.com/ldmonster/tts-parser/internal/module"

	"github.com/gabriel-vasile/mimetype"
	uberzap "go.uber.org/zap"
)

func NewClient(logger *uberzap.Logger) *Client {
	return &Client{
		client:                 &http.Client{},
		path:                   "tmp/",
		maxConcurrentDownloads: 3,
		errorCh:                make(chan error, 3),
		logger:                 logger,
	}
}

type Client struct {
	client *http.Client
	path   string

	maxConcurrentDownloads int

	errorCh chan error
	logger  *uberzap.Logger
}

func (c *Client) DownloadModule(ctx context.Context, mod *module.TTSModule) []service.File {
	downloadCh := make(chan struct{}, c.maxConcurrentDownloads)
	downloadedModules := make([]service.File, 0, 1)
	wg := new(sync.WaitGroup)

	// Sort module files by URL for consistent ordering
	files := c.getSortedModuleFiles(mod)

	// Start download worker goroutine
	go c.downloadWorker(ctx, files, downloadCh, wg, mod, &downloadedModules)

	// Wait for downloads and handle errors
	return c.waitForDownloads(ctx, downloadedModules)
}

func (c *Client) getSortedModuleFiles(mod *module.TTSModule) []module.ModuleFile {
	files := slices.Collect(maps.Values(mod.GetAll()))
	slices.SortFunc(files, func(a, b module.ModuleFile) int {
		return cmp.Compare(a.URL, b.URL)
	})
	return files
}

func (c *Client) downloadWorker(ctx context.Context, files []module.ModuleFile, downloadCh chan struct{}, wg *sync.WaitGroup, mod *module.TTSModule, downloadedModules *[]service.File) {
	defer close(c.errorCh)

	for _, mf := range files {
		if mf.URL == "" {
			continue
		}

		downloadCh <- struct{}{}
		wg.Add(1)

		go func(mf module.ModuleFile) {
			defer func() {
				<-downloadCh
				wg.Done()
			}()

			if c.fileExists(&mf) {
				c.logger.Info("file already exists", uberzap.String("url", mf.URL))
				return
			}

			if err := c.download(ctx, &mf); err != nil {
				c.errorCh <- fmt.Errorf("url: %s download: %w", mf.URL, err)
				return
			}

			c.logger.Info("downloaded", uberzap.String("url", mf.URL))

			*downloadedModules = append(*downloadedModules, service.File{
				ModuleID:  mod.ID,
				Type:      mf.Type,
				URL:       mf.URL,
				Extension: mf.GetExtension(),
			})
		}(mf)
	}

	wg.Wait()
}

func (c *Client) waitForDownloads(ctx context.Context, downloadedModules []service.File) []service.File {
	for {
		select {
		case err, ok := <-c.errorCh:
			if !ok {
				return downloadedModules
			}
			c.logger.Warn("download file", uberzap.Error(err))
		case <-ctx.Done():
			return downloadedModules
		}
	}
}

func (c *Client) fileExists(mf *module.ModuleFile) bool {
	abs, err := filepath.Abs(filepath.Join(c.path, mf.GetFolder(), mf.GetFilename()) + mf.GetExtension())
	if err != nil {
		c.logger.Warn("absolute path", uberzap.Error(err))
	}

	_, err = os.Stat(abs)
	return err == nil
}

var googleSignInRegex = regexp.MustCompile(`^accounts.google.com$`)

func (c *Client) download(ctx context.Context, mf *module.ModuleFile) error {
	// `https://steamusercontent-a.akamaihd.net/ugc/929306232365497323/03A7F5D6C7E7BC387121E8C444A9751CD81CCC9C/`
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, mf.URL, nil)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("user-agent", "curl/7.84.0")
	req.Header.Set("accept", "*/*")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	if googleSignInRegex.MatchString(resp.Request.URL.Host) {
		return errors.New("google 403")
	}

	body := io.Reader(resp.Body)

	if mf.GetExtension() == "" {
		mtype, recycledBody, err := detectMimeType(body)
		if err != nil {
			return fmt.Errorf("detecting mime type: %w", err)
		}
		body = recycledBody
		mf.Extension = mtype.Extension()
	}

	filename := filepath.Join(c.path, mf.GetFolder(), mf.GetFilename()) + mf.GetExtension()
	if err := saveFile(body, filename); err != nil {
		return fmt.Errorf("saving file: %w", err)
	}

	return nil
}

func saveFile(body io.Reader, filepath string) error {
	if err := os.MkdirAll(path.Dir(filepath), 0o777); err != nil {
		return fmt.Errorf("creating directories: %w", err)
	}

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0o777)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	if _, err := file.ReadFrom(body); err != nil {
		return fmt.Errorf("writing to file: %w", err)
	}

	return nil
}

func detectMimeType(input io.Reader) (*mimetype.MIME, io.Reader, error) {
	header := bytes.NewBuffer(nil)

	// After DetectReader, the data read from input is copied into header.
	mtype, err := mimetype.DetectReader(io.TeeReader(input, header))
	if err != nil {
		return nil, nil, fmt.Errorf("detecting mime type: %w", err)
	}

	// Concatenate back the header to the rest of the file.
	// recycled now contains the complete, original data.
	recycled := io.MultiReader(header, input)

	return mtype, recycled, nil
}
