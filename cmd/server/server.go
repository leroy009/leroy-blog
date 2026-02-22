package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leroy009/leroy-blog/internal/config"
	"github.com/leroy009/leroy-blog/internal/router"
)

func main() {
	cfg := config.Load()

	setupLogger(cfg.LogFile, cfg.LogLevel)

	r := router.New(cfg, slog.Default())

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		slog.Info("server starting", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	slog.Info("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited cleanly")
}

func setupLogger(logFile, logLevel string) {
	var level slog.Level
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}

	// Always write JSON to stderr — systemd/journald captures it automatically.
	// Use: journalctl -u leroy-blog -f | jq
	out := io.Writer(os.Stderr)

	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
		if err != nil {
			slog.Error("failed to open log file, falling back to stderr", "path", logFile, "error", err)
		} else {
			// Write to both stderr and the log file simultaneously.
			out = io.MultiWriter(os.Stderr, f)
		}
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(out, opts)))
}
