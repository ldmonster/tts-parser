package internal

import (
	"errors"

	"github.com/Masterminds/semver/v3"
)

var (
	ErrModuleConflict   = errors.New("module already exists")
	ErrModuleIsNotFound = errors.New("module is not found")
)

type Module struct {
	ID uint

	Name          string
	EpochTime     uint
	VersionNumber *semver.Version
}
