package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/v1tbrah/kvdb/config"
	"github.com/v1tbrah/kvdb/engine"
	"github.com/v1tbrah/kvdb/storage"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("parse config: %v", err)
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: cfg.LogLvl})))

	newStorage := storage.NewStorage()

	newEngine, err := engine.NewEngine(cfg.Server.Host, cfg.Server.Port, newStorage)
	if err != nil {
		slog.Error("new engine", "error", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	if err = newEngine.Launch(ctx); err != nil && err != context.Canceled {
		slog.Error("new engine. Launch", "error", err)
		cancel()
		os.Exit(1)
	}
	cancel()
	slog.Info("tcp server stopped")
}
