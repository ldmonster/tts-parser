package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ldmonster/tts-parser/internal/storage/gorm/model"
	"github.com/ldmonster/tts-parser/internal/storage/gorm/session"

	service "github.com/ldmonster/tts-parser/internal"

	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type File struct {
	DB *gorm.DB
}

func NewFile(db *gorm.DB) *File {
	return &File{
		DB: db,
	}
}

func (f *File) AutoMigrate(ctx context.Context) error {
	return session.DB(ctx, f.DB).Omit(clause.Associations).AutoMigrate(&model.File{})
}

func (f *File) Get(ctx context.Context, id uint) (*model.File, error) {
	existing := &model.File{
		ID: id,
	}

	db := session.DB(ctx, f.DB).Omit(clause.Associations).First(existing)

	return existing, db.Error
}

func (f *File) GetByTelegramID(ctx context.Context, telegramID int64) (*model.File, error) {
	existing := &model.File{}

	db := session.DB(ctx, f.DB).Omit(clause.Associations).Where("telegram_id = ?", telegramID).First(existing)
	if db.Error != nil && errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return nil, service.ErrFileIsNotFound
	}

	if db.Error != nil {
		return nil, db.Error
	}

	return existing, nil
}

func (f *File) ListByModuleID(ctx context.Context, id uint) ([]model.File, error) {
	existing := make([]model.File, 0, 1)

	db := session.DB(ctx, f.DB).Omit(clause.Associations).Where("module_id = ?", id).Find(&existing)

	if len(existing) > 0 {
		if existing[0].FileType == "" {
			panic("existing type are empty")
		}
	}

	return existing, db.Error
}

func (f *File) List(ctx context.Context) ([]model.File, error) {
	existing := make([]model.File, 0, 1)

	db := session.DB(ctx, f.DB).Omit(clause.Associations).Find(&existing)

	return existing, db.Error
}

func (f *File) Create(ctx context.Context, file *model.File) (*model.File, error) {
	db := session.DB(ctx, f.DB).Omit(clause.Associations).Create(file)
	sqliteErr := sqlite3.Error{}
	if db.Error != nil && errors.As(db.Error, &sqliteErr) && int(sqliteErr.ExtendedCode) == int(sqlite3.ErrConstraintUnique) {
		return nil, service.ErrFileConflict
	}

	if db.Error != nil {
		return nil, db.Error
	}

	return file, db.Error
}

func (f *File) BatchCreate(ctx context.Context, file ...model.File) error {
	if len(file) == 0 {
		return nil
	}

	db := session.DB(ctx, f.DB).Clauses(clause.OnConflict{DoNothing: true}).Create(file)
	sqliteErr := sqlite3.Error{}
	if db.Error != nil && errors.As(db.Error, &sqliteErr) && int(sqliteErr.ExtendedCode) == int(sqlite3.ErrConstraintUnique) {
		return service.ErrFileConflict
	}

	if db.Error != nil {
		return db.Error
	}

	return db.Error
}

func (f *File) Update(ctx context.Context, file *model.File) error {
	existing := &model.File{
		ID: file.ID,
	}

	db := session.DB(ctx, f.DB).Omit(clause.Associations).Where("id = ?", file.ID).First(existing)
	if db.Error != nil {
		return fmt.Errorf("first select: %w", db.Error)
	}

	db = session.DB(ctx, f.DB).Omit(clause.Associations).Model(existing).Updates(file)
	if db.Error != nil {
		return fmt.Errorf("update: %w", db.Error)
	}

	return nil
}

func (f *File) DeleteByModuleID(ctx context.Context, id uint) error {
	db := session.DB(ctx, f.DB).Omit(clause.Associations).Where("module_id = ?", id).Delete(&model.File{})
	if db.Error != nil {
		return fmt.Errorf("delete by module id: %w", db.Error)
	}

	return nil
}
