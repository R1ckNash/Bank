package config

import (
	"log"
	"os"
)

type Config struct {
	DBUrl     string
	Port      string
	JWTSecret string
}

func LoadConfig() *Config {
	cfg := &Config{
		DBUrl:     getEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/bank?sslmode=disable"),
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "supersecretkey"),
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("variable %s didn't set, using default value: %s", key, defaultValue)
		return defaultValue
	}
	return val
}
