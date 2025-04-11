package main

import (
	"flag"
	"log/slog"

	"github.com/av-ugolkov/bitly-copy/internal/app"
	"github.com/av-ugolkov/bitly-copy/internal/config"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "./configs/develop.yaml", "path to config file")

	flag.Parse()

	slog.Info("init config")
	cfg := config.Init(configPath)

	app.Init(cfg)
}
