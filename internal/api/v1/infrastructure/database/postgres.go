package database

import (
	"database/sql"

	"example.com/m/internal/api/v1/infrastructure/logger"
	"example.com/m/internal/config"
	"go.uber.org/zap"
)

var Db *sql.DB

func ConnectToDatabase() {
	db, err := sql.Open("postgres", config.Config.PostgresConnectionString)
	if err != nil {
		logger.Logger.Fatal("Failed connect to db: ", zap.Error(err))
	}

	err = db.Ping()
	if err != nil {
		logger.Logger.Fatal("Failed connect to db: ", zap.Error(err))
	}

	Db = db
}
