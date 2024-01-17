package config

import (
	"github.com/FaturFawkes/NOX-AI/pkg/utils"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func Environment() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			utils.Logger().Error("error while loading .env file", zap.Error(err))
		}
	} else {
		utils.Logger().Warn("running service without configuration from .env")
	}
}
