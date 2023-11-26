package config

import (
	"errors"
	"flag"
	"log/slog"
	"os"
)

// Config daemon
//
// Flags:
// h - host daemon
// p - port daemon
// l - log level (debug, info, warn, error)
// w - with writing to WAL
// s - sync writing to WAL
//
// Envs:
// HOST - host daemon
// PORT - port daemon
// LOG_LVL - log level (debug, info, warn, error)
// WAL - with writing to WAL
// SYNC_WAL - sync writing to WAL
type Config struct {
	Server
	LogLvl slog.Level
	// TODO add log format: console, json
	WithWritingToWAL bool
	SyncWritingToWAL bool
}

type Server struct {
	Port string
	Host string
}

func New() (Config, error) {
	host, port, logLvl := "localhost", "4321", "info"
	withWritingToWAL, syncWritingToWAL := true, true
	flag.StringVar(&host, "h", host, "host tcp server")
	flag.StringVar(&port, "p", port, "port tcp server")
	flag.StringVar(&logLvl, "l", logLvl, "log lvl")
	flag.BoolVar(&withWritingToWAL, "w", withWritingToWAL, "with writing to WAL")
	flag.BoolVar(&syncWritingToWAL, "s", syncWritingToWAL, "sync writing to WAL")

	flag.Parse()

	if envHost := os.Getenv("HOST"); envHost != "" {
		host = envHost
	}
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	if envLogLvl := os.Getenv("LOG_LVL"); envLogLvl != "" {
		logLvl = envLogLvl
	}
	if envWithWritingToWAL := os.Getenv("WAL"); envWithWritingToWAL != "" {
		withWritingToWAL = envWithWritingToWAL == "true"
	}
	if envSyncWritingToWAL := os.Getenv("SYNC_WAL"); envSyncWritingToWAL != "" {
		syncWritingToWAL = envSyncWritingToWAL == "true"
	}

	validLogLvls := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	if _, ok := validLogLvls[logLvl]; !ok {
		return Config{}, errors.New("invalid log lvl")
	}

	cfg := Config{}
	cfg.Server.Host = host
	cfg.Server.Port = port
	cfg.LogLvl = validLogLvls[logLvl]
	cfg.WithWritingToWAL = withWritingToWAL
	cfg.SyncWritingToWAL = syncWritingToWAL

	return cfg, nil
}
