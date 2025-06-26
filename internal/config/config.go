package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env             string
	GracefulTimeout time.Duration
	HTTPServer      HTTPServer
}

type HTTPServer struct {
	Protocol     string
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func MustLoad() *Config {
	return &Config{
		Env:             getEnv("APP_ENV", "local"),
		GracefulTimeout: getEnvAsDuration("GRACEFUL_TIMEOUT", 2*time.Second),
		HTTPServer: HTTPServer{
			Protocol:     getEnv("PROTOCOL", "http"),
			Host:         getEnv("HOST", "localhost"),
			Port:         getEnv("PORT", ":8080"),
			ReadTimeout:  getEnvAsDuration("READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getEnvAsDuration("WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getEnvAsDuration("IDLE_TIMEOUT", 120*time.Second),
		},
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
			log.Fatalf("Invalid duration for %s: %v. Using default: %v", key, err, defaultVal)
		}
	}
	return defaultVal
}
