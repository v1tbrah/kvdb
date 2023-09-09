package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/v1tbrah/kvdb/engine"
	"github.com/v1tbrah/kvdb/storage"
)

func main() {
	host, port, logLvl := "localhost", "4321", int(slog.LevelInfo)
	flag.StringVar(&host, "h", host, "host tcp server")
	flag.StringVar(&port, "p", port, "port tcp server")
	flag.IntVar(&logLvl, "l", logLvl, "log lvl")
	flag.Parse()

	if envHost := os.Getenv("HOST"); envHost != "" {
		host = envHost
	}
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	if envLogLvlStr := os.Getenv("LOGLVL"); envLogLvlStr != "" {
		envLogLvl, err := strconv.Atoi(envLogLvlStr)
		if err != nil {
			log.Fatalln("invalid log lvl")
		}
		logLvl = envLogLvl
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.Level(logLvl)})))

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
