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

	go func() {
		wg := new(sync.WaitGroup)

		result := slices.Collect(maps.Values(mod.GetAll()))
		slices.SortFunc(result, func(a, b module.ModuleFile) int {
			return cmp.Compare(a.URL, b.URL)
		})

		for _, mf := range result {
			if mf.URL == "" {
				continue
			}

			downloadCh <- struct{}{}
			wg.Add(1)

			go func() {
				defer func() {
					<-downloadCh
					wg.Done()
				}()

				if c.fileExists(&mf) {
					c.logger.Info("file already exists", uberzap.String("url", mf.URL))

					return
				}

				err := c.download(ctx, &mf)
				if err != nil {
					c.errorCh <- fmt.Errorf("url: %s download: %w", mf.URL, err)
					return
				}

				c.logger.Info("downloaded", uberzap.String("url", mf.URL))

				downloadedModules = append(downloadedModules, service.File{
					ModuleID:  mod.ID,
					Type:      mf.Type,
					URL:       mf.URL,
					Extension: mf.GetExtension(),
				})
			}()
		}

		wg.Wait()
		close(c.errorCh)
	}()

	for {
		select {
		case err, ok := <-c.errorCh:
			if ok {
				c.logger.Warn("download file", uberzap.Error(err))

				continue
			}

			return downloadedModules
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

	found := googleSignInRegex.MatchString(resp.Request.URL.Host)
	if found {
		return errors.New("google 403")
	}

	body := io.Reader(resp.Body)

	if mf.GetExtension() == "" {
		var mtype *mimetype.MIME

		mtype, body, err = recycleReader(body)
		if err != nil {
			return fmt.Errorf("recycle reader: %d", resp.StatusCode)
		}

		mf.Extension = mtype.Extension()
	}

	filename := filepath.Join(c.path, mf.GetFolder(), mf.GetFilename()) + mf.GetExtension()

	err = writeFile(body, filename)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func writeFile(body io.Reader, file string) error {
	err := os.MkdirAll(filepath.Dir(file), 0o777)
	if err != nil {
		return fmt.Errorf("mkdir all: %w", err)
	}

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0o777)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	_, err = f.ReadFrom(body)
	if err != nil {
		return fmt.Errorf("read from: %w", err)
	}

	return nil
}

// recycleReader returns the MIME type of input and a new reader
// containing the whole data from input.
func recycleReader(input io.Reader) (*mimetype.MIME, io.Reader, error) {
	// header will store the bytes mimetype uses for detection.
	header := bytes.NewBuffer(nil)

	// After DetectReader, the data read from input is copied into header.
	mtype, err := mimetype.DetectReader(io.TeeReader(input, header))
	if err != nil {
		return nil, nil, fmt.Errorf("detect reader: %w", err)
	}

	// Concatenate back the header to the rest of the file.
	// recycled now contains the complete, original data.
	recycled := io.MultiReader(header, input)

	return mtype, recycled, nil
}
