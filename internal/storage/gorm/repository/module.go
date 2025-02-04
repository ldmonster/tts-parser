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

type Module struct {
	DB *gorm.DB
}

func NewModule(db *gorm.DB) *Module {
	return &Module{
		DB: db,
	}
}

func (m *Module) AutoMigrate(ctx context.Context) error {
	return session.DB(ctx, m.DB).Omit(clause.Associations).AutoMigrate(&model.Module{})
}

func (m *Module) Get(ctx context.Context, id uint) (*model.Module, error) {
	existing := &model.Module{
		ID: id,
	}

	db := session.DB(ctx, m.DB).Omit(clause.Associations).First(existing)

	return existing, db.Error
}

func (m *Module) GetByTelegramID(ctx context.Context, telegramID int64) (*model.Module, error) {
	existing := &model.Module{}

	db := session.DB(ctx, m.DB).Omit(clause.Associations).Where("telegram_id = ?", telegramID).First(existing)
	if db.Error != nil && errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return nil, service.ErrModuleIsNotFound
	}

	if db.Error != nil {
		return nil, db.Error
	}

	return existing, nil
}

func (m *Module) List(ctx context.Context) ([]model.Module, error) {
	existing := make([]model.Module, 0, 1)

	db := session.DB(ctx, m.DB).Omit(clause.Associations).Find(&existing)

	return existing, db.Error
}

func (m *Module) Create(ctx context.Context, module *model.Module) (*model.Module, error) {
	db := session.DB(ctx, m.DB).Omit(clause.Associations).Create(module)
	sqliteErr := sqlite3.Error{}
	if db.Error != nil && errors.As(db.Error, &sqliteErr) && int(sqliteErr.ExtendedCode) == int(sqlite3.ErrConstraintUnique) {
		return nil, service.ErrModuleConflict
	}

	if db.Error != nil {
		return nil, db.Error
	}

	return module, db.Error
}

func (m *Module) BatchCreate(ctx context.Context, module ...model.Module) error {
	db := session.DB(ctx, m.DB).Clauses(clause.OnConflict{DoNothing: true}).Create(module)
	sqliteErr := sqlite3.Error{}
	if db.Error != nil && errors.As(db.Error, &sqliteErr) && int(sqliteErr.ExtendedCode) == int(sqlite3.ErrConstraintUnique) {
		return service.ErrModuleConflict
	}

	if db.Error != nil {
		return db.Error
	}

	return db.Error
}

func (m *Module) Update(ctx context.Context, module *model.Module) error {
	existing := &model.Module{
		ID: module.ID,
	}

	db := session.DB(ctx, m.DB).Omit(clause.Associations).Where("id = ?", module.ID).First(existing)
	if db.Error != nil {
		return fmt.Errorf("first select: %w", db.Error)
	}

	db = session.DB(ctx, m.DB).Omit(clause.Associations).Model(existing).Updates(module)
	if db.Error != nil {
		return fmt.Errorf("update: %w", db.Error)
	}

	return nil
}
