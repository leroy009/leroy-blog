package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func Load() Config {
	return Config{
		Port:         getEnv("PORT", "8080"),
		ReadTimeout:  getEnvDuration("READ_TIMEOUT", 5*time.Second),
		WriteTimeout: getEnvDuration("WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:  getEnvDuration("IDLE_TIMEOUT", 60*time.Second),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	seconds, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("invalid value for %s", key)
	}

	return time.Duration(seconds) * time.Second
}
