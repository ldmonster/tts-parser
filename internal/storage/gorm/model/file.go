package model

import (
	"database/sql/driver"

	service "github.com/ldmonster/tts-parser/internal"
)

type FileType string

const (
	FileTypeAssets FileType = "assets"
	FileTypeModel  FileType = "model"
	FileTypeImage  FileType = "image"
	FileTypePDF    FileType = "pdf"
	FileTypeAudio  FileType = "audio"
)

func getServiceFileTypes() []FileType {
	tps := make([]FileType, service.FileTypeOverall)

	tps[service.FileTypeAsset] = FileTypeAssets
	tps[service.FileTypeModel] = FileTypeModel
	tps[service.FileTypeImage] = FileTypeImage
	tps[service.FileTypePDF] = FileTypePDF
	tps[service.FileTypeAudio] = FileTypeAudio

	return tps
}

func getFileTypes() map[FileType]service.FileType {
	tps := map[FileType]service.FileType{
		FileTypeAssets: service.FileTypeAsset,
		FileTypeModel:  service.FileTypeModel,
		FileTypeImage:  service.FileTypeImage,
		FileTypePDF:    service.FileTypePDF,
		FileTypeAudio:  service.FileTypeAudio,
	}

	return tps
}

func remapFromServiceFileType(ft service.FileType) FileType {
	return getServiceFileTypes()[ft]
}

func remapToServiceFileType(ft FileType) service.FileType {
	return getFileTypes()[ft]
}

func (st *FileType) Scan(value interface{}) error {
	b, ok := value.(string)
	if !ok {
		panic("file type is not string")
	}

	*st = FileType(b)

	return nil
}

func (st FileType) Value() (driver.Value, error) {
	return string(st), nil
}

type File struct {
	ID uint `gorm:"primarykey"`

	ModuleID  uint     `gorm:"column:module_id"`
	FileType  FileType `gorm:"column:file_type;type:file_type;not null"`
	URL       string   `gorm:"unique;not null;column:url"`
	Extension string   `gorm:"column:extension"`
}

func RemapFromServiceFiles(input ...service.File) []File {
	result := make([]File, 0, len(input))

	for _, f := range input {
		result = append(result, *RemapFromServiceFile(&f))
	}

	return result
}

func RemapFromServiceFile(input *service.File) *File {
	return &File{
		ModuleID:  input.ModuleID,
		FileType:  remapFromServiceFileType(input.Type),
		URL:       input.URL,
		Extension: input.Extension,
	}
}

func RemapToServiceFiles(input ...File) []service.File {
	result := make([]service.File, 0, len(input))

	for _, f := range input {
		result = append(result, *RemapToServiceFile(&f))
	}

	return result
}

func RemapToServiceFile(input *File) *service.File {
	return &service.File{
		ModuleID:  input.ModuleID,
		Type:      remapToServiceFileType(input.FileType),
		URL:       input.URL,
		Extension: input.Extension,
	}
}
