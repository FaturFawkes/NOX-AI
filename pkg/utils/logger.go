package utils

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

func Logger() *zap.Logger {
	encodeConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encodeConfig, os.Stdout, zap.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	return logger
}
