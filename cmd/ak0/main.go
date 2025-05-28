package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/myshkins/ak0/internal/helpers"
	"github.com/myshkins/ak0/internal/logger"
	"github.com/myshkins/ak0/internal/metrics"
	"github.com/myshkins/ak0/internal/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

//go:embed dist/*
var dist embed.FS

func NewServerHandler(
	// config *Config
	clientRateLimiters *middleware.ClientRateLimiters,
	blockList *middleware.BlockList,
  staticFiles *embed.FS,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, blockList, clientRateLimiters, staticFiles)
	handler := otelhttp.NewHandler(mux, "/")

	return handler
}

func run(ctx context.Context, cfg helpers.Config,) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

  lp := cfg.LogPath
  logfile := logger.NewLogger(lp)
	go logger.ListenForLogrotate(lp, logfile, ctx)

	// set up otel
	otelShutdown, err := metrics.SetupOTelSDK(ctx, cfg)
	if err != nil {
		return err
	}

	defer func() {
		err = errors.Join(err, otelShutdown(ctx))
	}()

	crl := middleware.NewClientRateLimiters()
	bl := middleware.NewBlockList()
  srv := NewServerHandler(crl, bl, &dist)
	httpServer := &http.Server{
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
		// Addr: net.JoinHostPort(config.Host, config.Port),
		Addr:    net.JoinHostPort(cfg.ServerAddress, cfg.Port),
		Handler: srv,
	}

	go middleware.CleanupBlocklist(ctx, bl)
	go middleware.CleanupRateLimiters(ctx, crl)

	go func() {
		msg := fmt.Sprintf("listening on %v", httpServer.Addr)
		slog.Info(msg)
		if err := httpServer.ListenAndServe(); err != nil {
			msg := fmt.Sprintf("error listening and serving: %s", err)
			slog.Error(msg)
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown initiated")
	shutdownCtx := context.Background()
	shutdownCtx, shutDownCancel := context.WithTimeout(shutdownCtx, 10*time.Second)
	defer shutDownCancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		msg := fmt.Sprintf("error shutting down http server: %s", err)
		slog.Error(msg)
	}

	slog.Info("server shutdown complete")
	return nil
}

func main() {
	env := flag.String("env", "dev", "value to signal the type of environment to run in")
  configPath := flag.String("config-file", "./config.json", "path to the config file")
	flag.Parse()

  config, err := helpers.LoadConfig(*configPath, *env)
  if err != nil {
		log.Fatal(err)
	}

	err = os.Setenv("AK0_ENV", *env)
	if err != nil {
		log.Fatal(err)
	}

	if err := run(context.Background(), config); err != nil {
    fmt.Println(err.Error())
		os.Exit(1)
	}
}
