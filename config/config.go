package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_logger "github.com/jesslyn-ctrl/doit-url-shortener/pkg/logger"
	"github.com/joho/godotenv"
)

var logs = _logger.GetContextLoggerf(nil)

// Config structure for the application
type Config struct {
	App AppConfig
	URL URLConfig
}

// AppConfig holds app related settings
type AppConfig struct {
	Name string `env:"APP_NAME"`
	Port string `env:"APP_PORT"`
}

// URLConfig holds url shortener related settings
type URLConfig struct {
	UrlTTLSeconds time.Duration `env:"URL_TTL_SECONDS"`
}

// AppConfigInstance Global variable to store the loaded config
var AppConfigInstance Config

// LoadConfig reads from .env file and sets values in AppConfigInstance
func LoadConfig() error {
	currentWorkDirectory, _ := os.Getwd()
	loadEnvFile(currentWorkDirectory + "/.env")

	urlCfg, err := loadUrlConfig()
	if err != nil {
		return err
	}

	AppConfigInstance = Config{
		App: AppConfig{
			Name: os.Getenv("APP_NAME"),
			Port: os.Getenv("APP_PORT"),
		},
		URL: urlCfg,
	}

	logs.Info("✅ Config loaded successfully!")
	return nil
}

func loadEnvFile(filepath string) {
	if err := godotenv.Load(filepath); err != nil {
		logs.Warnf("Warning: No %s file found or error loading: %v", filepath, err)
	}
}

func loadUrlConfig() (URLConfig, error) {
	ttlStr := os.Getenv("URL_TTL_SECONDS")
	if ttlStr == "" {
		ttlStr = "86400" // 24 hours
	}

	ttlSeconds, err := strconv.Atoi(ttlStr)
	if err != nil || ttlSeconds <= 0 {
		return URLConfig{}, fmt.Errorf("invalid URL_TTL_SECONDS: %s", ttlStr)
	}

	return URLConfig{
		UrlTTLSeconds: time.Duration(ttlSeconds) * time.Second,
	}, nil
}
