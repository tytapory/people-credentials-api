package config

import (
	"github.com/joho/godotenv"
	"os"
	"people-credentials-api/pkg/logger"
	"sync"
)

var (
	once sync.Once
	cfg  *Config
)

// Config - структура для хранения конфигурации сервиса
type Config struct {
	ServerPort      string
	DatabasePort    string
	DatabaseUser    string
	DatabasePass    string
	DatabaseName    string
	DatabaseHost    string
	DatabaseSSLMode string
	LogLevel        string
}

// Get загружает конфигурацию из переменных окружения (только при первом вызове)
// и возвращает указатель на структуру Config
func Get() *Config {
	logger.Debug("Get() called")

	once.Do(func() {
		logger.Info("Loading configuration for the first time")

		err := godotenv.Load()
		if err != nil {
			logger.Warn("Could not load .env file, proceeding with environment variables only")
		} else {
			logger.Info(".env file successfully loaded")
		}

		cfg = &Config{
			ServerPort:      getEnv("PEOPLE_CREDENTIALS_SERVER_PORT", "8080", os.LookupEnv),
			DatabasePort:    getEnv("PEOPLE_CREDENTIALS_DATABASE_PORT", "5432", os.LookupEnv),
			DatabaseUser:    getEnv("PEOPLE_CREDENTIALS_DATABASE_USER", "postgres", os.LookupEnv),
			DatabasePass:    getEnv("PEOPLE_CREDENTIALS_DATABASE_PASSWORD", "password", os.LookupEnv),
			DatabaseName:    getEnv("PEOPLE_CREDENTIALS_DATABASE_NAME", "user_creds_db", os.LookupEnv),
			DatabaseHost:    getEnv("PEOPLE_CREDENTIALS_DATABASE_HOST", "localhost", os.LookupEnv),
			DatabaseSSLMode: getEnv("PEOPLE_CREDENTIALS_DATABASE_SSL_MODE", "disable", os.LookupEnv),
			LogLevel:        getEnv("PEOPLE_CREDENTIALS_LOG_LEVEL", "info", os.LookupEnv),
		}

		logger.Info("Configuration successfully loaded and cached")
	})

	logger.Debug("Returning cached configuration")
	return cfg
}

// getEnv получает значение переменной окружения по ключу.
// Если переменная не задана, возвращает значение по умолчанию.
func getEnv(key, fallback string, getEnvFunc func(string) (string, bool)) string {
	logger.Debug("Trying to load environment variable: " + key)

	if value, ok := getEnvFunc(key); ok {
		logger.Info("Loaded environment variable: " + key + " = " + value)
		return value
	}

	logger.Warn("Environment variable not found: " + key + ", using fallback: " + fallback)
	return fallback
}
