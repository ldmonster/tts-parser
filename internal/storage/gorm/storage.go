package gorm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ldmonster/tts-parser/internal/storage/gorm/repository"
	"github.com/ldmonster/tts-parser/internal/storage/gorm/session"

	"github.com/ldmonster/tts-parser/internal/storage/gorm/model"

	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"gorm.io/gorm/logger"

	// "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	service "github.com/ldmonster/tts-parser/internal"

	uberzap "go.uber.org/zap"
	"gorm.io/gorm"
)

type Storage struct {
	gorm   session.Gorm
	Module *repository.Module
	File   *repository.File

	logger *uberzap.Logger
}

func NewStorage(dbpath string, l *uberzap.Logger) (*Storage, error) {
	dir := filepath.Dir(dbpath)

	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return nil, fmt.Errorf("mkdir to db: %w", err)
	}

	gcfg := &gorm.Config{}

	if l.Level() == uberzap.DebugLevel {
		gcfg.Logger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{LogLevel: logger.Info})
	}

	db, err := gorm.Open(sqlite.Open(dbpath), gcfg)
	if err != nil {
		return nil, fmt.Errorf("open db connection: %w", err)
	}

	return &Storage{
		gorm:   session.GORM(db, &sql.TxOptions{}),
		Module: repository.NewModule(db),
		File:   repository.NewFile(db),
		logger: l,
	}, nil
}

func (s *Storage) AutoMigrate(ctx context.Context) error {
	session, err := s.gorm.Begin(ctx)
	if err != nil {
		return err
	}
	defer session.Rollback()

	err = s.Module.AutoMigrate(ctx)
	if err != nil {
		return err
	}

	err = s.File.AutoMigrate(ctx)
	if err != nil {
		return err
	}

	err = session.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) StartTransaction(ctx context.Context) (session.Session, error) {
	return s.gorm.Begin(ctx)
}

func (s *Storage) Register(ctx context.Context, telegramID int64, modulename string, lastMessageID int, startingBalance map[string]int64) error {
	newModule := &model.Module{}

	newModule, cErr := s.Module.Create(ctx, newModule)
	if cErr != nil && errors.Is(cErr, service.ErrModuleConflict) {
		s.logger.Error("creating module", uberzap.Any("module", newModule), uberzap.Error(cErr))
		return fmt.Errorf("creating module: %w", cErr)
	}

	return nil
}

func (s *Storage) UpdateModule(ctx context.Context, telegramID int64, fn func(u *model.Module) error) error {
	module, err := s.Module.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return fmt.Errorf("get by telegram id: %w", err)
	}

	err = fn(module)
	if err != nil {
		return fmt.Errorf("mutating func: %w", err)
	}

	err = s.Module.Update(ctx, module)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}
