package main

import (
	"fmt"
	"log"
	"net/http"

	_config "github.com/jesslyn-ctrl/doit-url-shortener/config"
	_http "github.com/jesslyn-ctrl/doit-url-shortener/internal/http"
	_logger "github.com/jesslyn-ctrl/doit-url-shortener/pkg/logger"
)

var logs = _logger.GetContextLoggerf(nil)

func main() {
	// 1. Load config from .env
	if err := _config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	cfg := _config.AppConfigInstance

	// 2. Init HTTP router
	router := _http.NewRouter()

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	logs.Infof("starting server at %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
