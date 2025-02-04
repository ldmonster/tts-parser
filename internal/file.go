package internal

import (
	"errors"
)

var (
	ErrFileConflict   = errors.New("file already exists")
	ErrFileIsNotFound = errors.New("file is not found")
)

type FileType int

const (
	FileTypeAsset = iota
	FileTypeModel
	FileTypeImage
	FileTypePDF
	FileTypeAudio
	FileTypeOverall
)

type File struct {
	ID uint

	ModuleID  uint
	Type      FileType
	URL       string
	Extension string
}
