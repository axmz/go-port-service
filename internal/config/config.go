package config

import "os"

type Config struct {
	Port     string
	Env      string
	DBURL    string
	LogLevel string
}

func LoadConfig() *Config {
	return &Config{
		Port:     getEnv("PORT", ":8080"),
		Env:      getEnv("APP_ENV", "development"),
		DBURL:    getEnv("DATABASE_URL", "postgres://localhost:5432/mydb"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
