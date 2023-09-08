package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/v1tbrah/kvdb/engine"
	"github.com/v1tbrah/kvdb/storage"
)

func main() {
	host, port := "localhost", "4321"
	flag.StringVar(&host, "h", host, "host tcp server")
	flag.StringVar(&port, "p", port, "port tcp server")

	if envHost := os.Getenv("HOST"); envHost != "" {
		host = envHost
	}
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	newStorage := storage.NewStorage()

	newEngine, err := engine.NewEngine(host, port, newStorage)
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
