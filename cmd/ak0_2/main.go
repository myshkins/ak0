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
	"sync"
	"time"

	"github.com/myshkins/ak0_2/internal/logger"
	"github.com/myshkins/ak0_2/internal/metrics"
	"github.com/myshkins/ak0_2/internal/middleware"
  "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)


var meter = otel.Meter("github.com/myshkins/ak0_2")

func testCounter() {
	apiCounter, err := meter.Int64Counter(
		"test.counter",
		metric.WithDescription("Number of API calls."),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		panic(err)
	}
  ctx := context.Background()
  c := context.WithValue(ctx, "ak02contextest", "hewo")
  fmt.Println(c.Value("ak02contextest"))
  apiCounter.Add(c, 1)
}

func NewServerHandler(
  // config *Config
  // commentStore *commentStore
  // BotList
  clientRateLimiters *middleware.ClientRateLimiters,
) http.Handler {
  mux := http.NewServeMux()
  addRoutes(mux, clientRateLimiters)
  handler := otelhttp.NewHandler(mux, "/")

  go middleware.CleanupRateLimiters(clientRateLimiters)
  // var handler http.Handler = mux
  return handler
}

func run(ctx context.Context, w io.Writer, slogger *slog.Logger, args []string) error {
  testCounter()
  ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
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
  srv := NewServerHandler(crl)
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

  var wg sync.WaitGroup
  wg.Add(1)
  go func() {
    defer wg.Done()
    <-ctx.Done()
    shutdownCtx := context.Background()
    shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10 * time.Second)
    defer cancel()
    if err := httpServer.Shutdown(shutdownCtx); err != nil {
      msg := fmt.Sprintf("error shutting down http server: %s", err)
      slog.Error(msg)
    }
  }()
  wg.Wait()

  return nil
}


func main() {
  env := flag.String("env", "dev", "value to signal the type of environment to run in")
  flag.Parse()

  err := os.Setenv("AK0_2_ENV", *env)
  if err != nil {
    log.Fatal(err)
  }

  slogger := logger.NewLogger()
  slog.SetDefault(slogger)

  ctx := context.Background()
  if err := run(ctx, os.Stdout, slogger, os.Args); err != nil {
    slogger.Error(err.Error())
    os.Exit(1)
  }
}

