package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
  "syscall"
	"time"

	"github.com/myshkins/ak0/internal/logger"
	"github.com/myshkins/ak0/internal/metrics"
	"github.com/myshkins/ak0/internal/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)


func NewServerHandler(
  // config *Config
  clientRateLimiters *middleware.ClientRateLimiters,
  blockList *middleware.BlockList,
) http.Handler {
  mux := http.NewServeMux()
  addRoutes(mux, blockList, clientRateLimiters)
  handler := otelhttp.NewHandler(mux, "/")

  go middleware.CleanupRateLimiters(clientRateLimiters) // move this to run?
  return handler
}

func run(ctx context.Context, w io.Writer, slogger *slog.Logger, args []string) error {
  ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
  defer cancel()

  // set up otel
  otelShutdown, err := metrics.SetupOTelSDK(ctx)
  if err != nil {
    return err
  }

  defer func() {
    err = errors.Join(err, otelShutdown(context.Background()))
  }()
  
  crl := middleware.NewClientRateLimiters()
  bl := middleware.NewBlockList()
  srv := NewServerHandler(crl, bl)
  httpServer := &http.Server{
    ReadTimeout: 120 * time.Second,
    WriteTimeout: 120 * time.Second,
    IdleTimeout: 120 * time.Second,
    // Addr: net.JoinHostPort(config.Host, config.Port),
    Addr: net.JoinHostPort("0.0.0.0", "8200"),
    Handler: srv,
  }

  go func() {
    msg := fmt.Sprintf("listening on %v\n", httpServer.Addr)
    slog.Info(msg)
    if err := httpServer.ListenAndServe(); err != nil {
      msg := fmt.Sprintf("error listening and serving: %s", err)
      slog.Error(msg)
    }
  }()

  <-ctx.Done()
  slogger.Info("shutdown initiated")
  shutdownCtx := context.Background()
  shutdownCtx, shutDownCancel := context.WithTimeout(shutdownCtx, 10 * time.Second)
  defer shutDownCancel()
  if err := httpServer.Shutdown(shutdownCtx); err != nil {
    msg := fmt.Sprintf("error shutting down http server: %s", err)
    slog.Error(msg)
  }

  slogger.Info("server shutdown complete")
  return nil
}


func main() {
  env := flag.String("env", "dev", "value to signal the type of environment to run in")
  flag.Parse()

  err := os.Setenv("AK0_ENV", *env)
  if err != nil {
    log.Fatal(err)
  }

  slogger, logfile := logger.NewLogger()
  logger.ListenForLogrotate(logfile)

  ctx := context.Background()
  if err := run(ctx, os.Stdout, slogger, os.Args); err != nil {
    slogger.Error(err.Error())
    os.Exit(1)
  }
}

