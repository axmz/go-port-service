package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port         string
	Env          string
	DBURL        string
	LogLevel     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func LoadConfig() *Config {
	return &Config{
		Port:         getEnv("PORT", ":8080"),
		Env:          getEnv("APP_ENV", "development"),
		DBURL:        getEnv("DATABASE_URL", "postgres://localhost:5432/mydb"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		ReadTimeout:  getEnvAsDuration("READ_TIMEOUT", 5*time.Second),
		WriteTimeout: getEnvAsDuration("WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:  getEnvAsDuration("IDLE_TIMEOUT", 120*time.Second),
	}
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if valStr, ok := os.LookupEnv(key); ok {
		if valInt, err := strconv.Atoi(valStr); err == nil {
			return time.Duration(valInt) * time.Second
		} else {
			log.Printf("Invalid duration for %s: %v. Using default: %v", key, err, defaultVal)
		}
	}
	return defaultVal
}
