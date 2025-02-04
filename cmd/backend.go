package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
	"time"
	"tts/internal/downloader"
	"tts/internal/module"
	"tts/internal/storage/gorm"
	"tts/internal/storage/gorm/model"

	uberzap "go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

type backend struct {
	cfg    *Config
	logger *uberzap.Logger

	storage *gorm.Storage

	bot *tele.Bot
}

func NewBackend(cfg *Config, logger *uberzap.Logger) *backend {
	return &backend{
		cfg:    cfg,
		logger: logger,
	}
}

func (be *backend) init() error {
	err := be.initStorage()
	if err != nil {
		return fmt.Errorf("storage initialization: %w", err)
	}

	return nil
}

func (be *backend) initStorage() error {
	ctx := context.Background()

	var err error
	dbPath := filepath.Join(be.cfg.RootPath, be.cfg.Storage.SqliteDBPath)
	be.storage, err = gorm.NewStorage(dbPath, be.logger.Named("storage"))
	if err != nil {
		return fmt.Errorf("creating storage: %w", err)
	}

	err = be.storage.AutoMigrate(ctx)
	if err != nil {
		return fmt.Errorf("auto migration: %w", err)
	}

	return nil
}

var jsonRegex = regexp.MustCompile(`^([0-9]*).json$`)

func (be *backend) Start(ctx context.Context) {
	dir := `E:\Tabletop Simulator\Mods\Workshop`

	fs, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	parsingWg := new(sync.WaitGroup)

	trottleCh := make(chan struct{}, 10)
	modulesCh := make(chan module.TTSModule, 100)

	dbWritingIsDoneCh := make(chan struct{}, 1)

	go func() {
		for {
			select {
			case mod, ok := <-modulesCh:
				{
					if !ok {
						dbWritingIsDoneCh <- struct{}{}
					}

					files, err := be.storage.File.ListByModuleID(ctx, mod.ID)
					if err != nil {
						panic(err)
					}

					orphans := mod.MergeFiles(model.RemapToServiceFiles(files...))
					if len(orphans) > 0 {
						be.logger.Warn("orphans", uberzap.Any("orphans", orphans))
						err := be.storage.File.DeleteByModuleID(ctx, orphans[0].ModuleID)
						if err != nil {
							panic(err)
						}
					}

					c := downloader.NewClient(be.logger)
					dm := c.DownloadModule(ctx, &mod)

					// TODO: count errors (origin len - downloaded len)
					be.logger.Info("module downloaded", uberzap.String("module", mod.Name), uberzap.Uint("id", mod.ID), uberzap.Int("count", len(dm)))

					err = be.storage.File.BatchCreate(ctx, model.RemapFromServiceFiles(dm...)...)
					if err != nil {
						panic(err)
					}

					be.storage.Module.Create(ctx, &model.Module{
						ID:            mod.ID,
						Name:          mod.Name,
						EpochTime:     mod.EpochTime,
						VersionNumber: mod.VersionNumber.Original(),
					})
				}
			case <-ctx.Done():
				fmt.Println("stopped")
				os.Exit(0)
			}
		}
	}()

	for _, f := range fs {
		parsingWg.Add(1)
		trottleCh <- struct{}{}

		go func() {
			defer func() {
				parsingWg.Done()
				<-trottleCh
			}()

			subs := jsonRegex.FindAllStringSubmatch(f.Name(), 1)

			if len(subs) == 0 {
				return
			}

			id, err := strconv.Atoi(subs[0][1])
			if err != nil {
				panic(err)
			}

			f, err := os.OpenFile(filepath.Join(dir, f.Name()), os.O_RDONLY, 0o666)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			mod := new(module.Module)
			err = json.NewDecoder(f).Decode(mod)
			if err != nil {
				panic(err)
			}

			timestamp, err := time.Parse("1/2/2006 15:04:05 PM", mod.Date)
			if err != nil {
				timestamp, err = time.Parse("01/02/2006 15:04:05", mod.Date)
				if err != nil {
					panic(err)
				}
			}

			result := module.NewTTSModule()
			result.ScanModule(mod)

			result.ID = uint(id)
			result.EpochTime = uint(timestamp.Unix())

			modulesCh <- *result
		}()
	}

	parsingWg.Wait()

	fmt.Println("parsing is done")

	close(modulesCh)

	<-dbWritingIsDoneCh
}
