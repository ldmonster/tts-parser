package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	env "github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	uberzap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

type StorageConfig struct {
	SqliteDBPath string `env:"SQLITE_DB_PATH" envDefault:"tmp/gorm.db"`
}

func newStorageConfiig() *StorageConfig {
	return &StorageConfig{}
}

type Config struct {
	Storage *StorageConfig `envPrefix:"STORAGE_"`

	RootPath string

	ConfigPath string `env:"CONFIG_PATH"`

	LogLevelRaw string              `env:"LOG_LEVEL" envDefault:"INFO"`
	LogLevel    uberzap.AtomicLevel `env:"-"`

	ZapLogLevel zapcore.Level
}

func NewConfig() *Config {
	return &Config{
		Storage:     newStorageConfiig(),
		ZapLogLevel: zapcore.ErrorLevel,
	}
}

var ErrEnvFileIsNotFound = errors.New(".env file is not found")
var ErrConfigFileIsNotFound = errors.New(".config.yaml file is not found")

func (cfg *Config) AutoLoadEnvs() error {
	absBinaryPath := filepath.Dir(os.Args[0])
	absEnvPath := filepath.Join(absBinaryPath, ".env")
	absConfigPath := filepath.Join(absBinaryPath, ".config.yaml")

	var envPath, configPath string

	if _, err := os.Stat(absEnvPath); err == nil {
		envPath = absEnvPath
	}

	if _, err := os.Stat(absConfigPath); err == nil {
		configPath = absConfigPath
	}

	if _, err := os.Stat("./.env"); err == nil {
		envPath = "./.env"
	}

	if _, err := os.Stat("./.config.yaml"); err == nil {
		configPath = "./.config.yaml"
	}

	if envPath == "" {
		return ErrEnvFileIsNotFound
	}

	if configPath == "" {
		return ErrConfigFileIsNotFound
	}

	cfg.ConfigPath = configPath

	cfg.RootPath = filepath.Dir(configPath)

	return godotenv.Load(envPath)
}

func (cfg *Config) Parse() error {
	opts := env.Options{
		Prefix: "",
	}

	err := env.ParseWithOptions(cfg, opts)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	cfg.LogLevel, err = uberzap.ParseAtomicLevel(cfg.LogLevelRaw)
	if err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
	}

	_, err = os.Stat(cfg.ConfigPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("config file info: %w", err)
	}

	f, err := os.Open(cfg.ConfigPath)
	if err != nil {
		return fmt.Errorf("config file is not open: %w", err)
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(cfg)
	if err != nil {
		return fmt.Errorf("config file decode: %w", err)
	}

	return nil
}
