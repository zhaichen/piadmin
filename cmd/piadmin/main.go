package main

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/ptmind/piadmin/internal/api"
	"github.com/ptmind/piadmin/internal/auth"
	"github.com/ptmind/piadmin/internal/config"
	"github.com/ptmind/piadmin/internal/monitor"
	"github.com/ptmind/piadmin/internal/server"
	"github.com/ptmind/piadmin/web"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	configPath := flag.String("config", "", "config file path")
	showVersion := flag.Bool("version", false, "show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("piadmin %s (built %s)\n", version, buildTime)
		return
	}

	cfg := config.Load(*configPath)

	collector := monitor.NewCollector(cfg.MonitorInterval())
	collector.Start()
	defer collector.Stop()

	a := auth.New(cfg.Auth)

	assets, err := web.Assets()
	if err != nil {
		slog.Warn("frontend assets not found, running API only", "err", err)
	}

	router := api.NewRouter(cfg, a, collector, assets)

	slog.Info("piadmin starting", "version", version, "addr", cfg.Server.Addr)
	server.Run(cfg.Server.Addr, router)
}
