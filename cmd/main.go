package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	_config "github.com/jesslyn-ctrl/doit-url-shortener/config"
	_domainUrl "github.com/jesslyn-ctrl/doit-url-shortener/internal/domain/url"
	_http "github.com/jesslyn-ctrl/doit-url-shortener/internal/http"
	_storage "github.com/jesslyn-ctrl/doit-url-shortener/internal/storage"
	_pkgClk "github.com/jesslyn-ctrl/doit-url-shortener/pkg/clock"
	_logger "github.com/jesslyn-ctrl/doit-url-shortener/pkg/logger"
)

var logs = _logger.GetContextLoggerf(context.TODO())

func main() {
	// 1. Load config from .env
	if err := _config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	cfg := _config.AppConfigInstance

	// 2. Init infra
	store := _storage.NewMemoryStore()
	clk := _pkgClk.NewSystemClock()

	// 3. Init domain service
	urlSvc := _domainUrl.NewService(store, clk)

	// 4. Init HTTP router
	router := _http.NewRouter(
		urlSvc,
		cfg.URL.UrlTTLSeconds,
	)

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	logs.Infof("starting server at %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
