package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/v1tbrah/kvdb/config"
	"github.com/v1tbrah/kvdb/dbengine"
	"github.com/v1tbrah/kvdb/memory"
	"github.com/v1tbrah/kvdb/server"
	"github.com/v1tbrah/kvdb/wal"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("config.New", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: cfg.LogLvl})))

	newWAL, err := wal.New(false)
	if err != nil {
		slog.Error("wal.New", slog.String("error", err.Error()))
		os.Exit(1)
	}

	newDBEngine, err := dbengine.New(memory.New[string, string](), newWAL)
	if err != nil {
		slog.Error("dbengine.New", slog.String("error", err.Error()))
		os.Exit(1)
	}

	newServer, err := server.New(cfg.Server.Host, cfg.Server.Port, newDBEngine)
	if err != nil {
		slog.Error("server.New", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	if err = newServer.Launch(ctx); err != nil && !errors.Is(err, context.Canceled) {
		cancel()
		slog.Error("newServer.Launch", slog.String("error", err.Error()))
		os.Exit(1)
	}
	cancel()

	if err = newWAL.Close(ctx); err != nil {
		slog.Error("newWAL.Close", slog.String("error", err.Error()))
	}

	slog.Info("tcp server stopped")
}
