package configs

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
)

type RedisConfig struct {
	Addr     string `envconfig:"REDIS_ADDR" default:"localhost:6379"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
	Username string `envconfig:"REDIS_USERNAME" default:""`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
}

type RepoConfig struct {
	Redis   RedisConfig
	Timeout Duration `envconfig:"REPO_REDIS_TIMEOUT" default:"5s"`
}

type BasicService struct {
	Redis RepoConfig
}

type ServicesConfig struct {
	Basic      RepoConfig
	Monitoring PollMonitoringService
}

type WebSocket struct {
	Port string `envconfig:"WEBSOCKET_PORT"`
}

type PollMonitoringService struct {
	WebSocket WebSocket
}

type AppConfig struct {
	Repo RepoConfig
	Srv  ServicesConfig
}

func LoadConfig() (*AppConfig, error) {
	var config AppConfig

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: .env file not loaded, continuing with environment variables only")
		}
	}

	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
		return nil, err
	}

	return &config, nil
}
