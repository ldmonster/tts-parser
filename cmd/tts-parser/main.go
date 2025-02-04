package main

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"syscall"
	"tts/internal/zap"

	uberzap "go.uber.org/zap"
)

func main() {
	cfg := NewConfig()

	err := cfg.AutoLoadEnvs()
	if err != nil && (!errors.Is(err, ErrEnvFileIsNotFound) && !errors.Is(err, ErrConfigFileIsNotFound)) {
		panic(err)
	}

	if err != nil {
		fmt.Println(err)
	}

	err = cfg.Parse()
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewProductionZaplogger("log.txt", cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	logger = logger.Named("tts")

	defer func(logger *uberzap.Logger) {
		_ = logger.Sync()
	}(logger)

	b := NewBackend(cfg, logger)

	err = b.init()
	if err != nil {
		logger.Fatal("backend initialization", uberzap.Error(err))
	}
	logger.Error("download module", uberzap.Error(err))

	ctx, stopNotify := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	b.Start(ctx)

	stopNotify()
}
