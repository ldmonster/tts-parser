package zap

import (
	"fmt"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewFileLogger(filename string, lvl zap.AtomicLevel) (*zap.Logger, error) {
	// Configure the time format
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create file and console encoders
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	// Open the log file
	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	// Create writers for file and console
	fileWriter := zapcore.AddSync(logFile)
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Set the log level
	defaultLogLevel := lvl

	// Create cores for writing to the file and console
	fileCore := zapcore.NewCore(fileEncoder, fileWriter, defaultLogLevel)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, defaultLogLevel)

	// Combine cores
	core := zapcore.NewTee(fileCore, consoleCore)

	// Create the logger with additional context information (caller, stack trace)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger, nil
}

// NewProductionZapLogger will return a new production logger backed by zap
func NewProductionZaplogger(filename string, lvl zap.AtomicLevel) (*zap.Logger, error) {
	conf := zap.NewProductionConfig()
	conf.Level = lvl

	conf.DisableCaller = true
	conf.DisableStacktrace = true

	zapLogger, err := conf.Build(zap.WrapCore(zapCore))
	if err != nil {
		return nil, fmt.Errorf("build: %w", err)
	}

	return zapLogger, nil
}

func zapCore(c zapcore.Core) zapcore.Core {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/chat.log",
		MaxSize:    2, // megabytes
		MaxBackups: 30,
		MaxAge:     28, // days
	})

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.DebugLevel,
	)
	cores := zapcore.NewTee(c, core)

	return cores
}
