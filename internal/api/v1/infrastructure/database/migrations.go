package database

import (
	"example.com/m/internal/api/v1/infrastructure/logger"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

func MigrateDB() {
	goose.SetBaseFS(nil)
	if err := goose.SetDialect("postgres"); err != nil {
		logger.Logger.Fatal("Failed connect to db: ", zap.Error(err))
	}

	if err := goose.Up(Db, "/scripts"); err != nil {
		logger.Logger.Fatal("Failed connect to db: ", zap.Error(err))
	}
}
