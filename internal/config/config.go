package config

import (
	"fmt"
	"os"
	"time"
)

var (
	Config   AppConfig = AppConfig{}
	IAMToken string
)

type AppConfig struct {
	PostgresConnectionString  string
	JWTSecret                 string
	JWTExpiration             time.Duration
	RedisConnectionString     string
	RedisPassword             string
	S3Region                  string
	S3AWSAccessKeyID          string
	S3AwsSecretAccessKey      string
	S3Endpoint                string
	YandexOAuthToken          string
	YandexCatalogID           string
	FirebasePathToCredentials string
	ChatBotPrompt             string
}

func InitConfig() {
	Config = AppConfig{
		PostgresConnectionString: fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME")),
		JWTSecret:                 os.Getenv("JWT_SECRET"),
		JWTExpiration:             time.Hour * 24 * 90,
		RedisConnectionString:     os.Getenv("REDIS_CONNECTION"),
		RedisPassword:             "",
		S3Region:                  os.Getenv("S3_REGION"),
		S3AWSAccessKeyID:          os.Getenv("S3_AWS_ACCESS_KEY_ID"),
		S3AwsSecretAccessKey:      os.Getenv("S3_AWS_SECRET_ACCESS_KEY"),
		S3Endpoint:                os.Getenv("S3_ENDPOINT"),
		YandexOAuthToken:          os.Getenv("YANDEX_OAUTH_TOKEN"),
		YandexCatalogID:           os.Getenv("YANDEX_CATALOG_ID"),
		FirebasePathToCredentials: os.Getenv("FIREBASE_PATH_TO_CREDENTIALS"),
		ChatBotPrompt:             os.Getenv("CHAT_BOT_PROMPT"),
	}
}
