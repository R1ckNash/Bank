package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Env       string `yaml:"env" env-default:"local" env-required:"true"`
	DBUrl     string `yaml:"dbUrl" env-required:"true"`
	Port      int    `yaml:"port" env-default:"8081"`
	JWTSecret string `yaml:"jwt-secret" env-default:"supersecretkey"`

	AuthService struct {
		Host string `yaml:"host" env-default:"bank-auth-service"`
		Port int    `yaml:"port" env-default:"8080"`
	} `yaml:"auth_service"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Config path could not be empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exists: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("could not read config: %s", err)
	}

	return &cfg
}
