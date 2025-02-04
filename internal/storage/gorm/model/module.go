package model

type Module struct {
	ID uint `gorm:"primarykey"`

	Name          string
	EpochTime     uint
	VersionNumber string
	// TelegramID       int64 `gorm:"unique;not null"`
}
