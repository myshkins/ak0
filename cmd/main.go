package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
  "log"
)

var logger = log.New(os.Stdout, "log_ak0_2: ", log.Ldate)

func NewServer(
  logger *log.Logger,
  // config *Config
  // commentStore *commentStore
  // anotherStore *anotherStore
) http.Handler {
  mux := http.NewServeMux()
  addRoutes(mux, logger)

  var handler http.Handler = mux
  // handler = someMiddleware(handler)
  return handler
}

func run(ctx context.Context, w io.Writer, args []string) error {
  ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
  defer cancel()
  
  srv := NewServer(logger)
  httpServer := &http.Server{
    // Addr: net.JoinHostPort(config.Host, config.Port),
    Addr: net.JoinHostPort("127.0.0.1", "8200"),
    Handler: srv,
  }

  go func() {
    log.Printf("listening on %s\n", httpServer.Addr)
    if err := httpServer.ListenAndServe(); err != nil {
      fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
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
      fmt.Fprint(os.Stderr, "error shuttind down http server: %s\n", err)
    }
  }()
  wg.Wait()

  return nil
}

func main() {
  ctx := context.Background()
  if err := run(ctx, os.Stdout, os.Args); err != nil {
    fmt.Fprintf(os.Stderr, "%s\n", err)
    os.Exit(1)
  }
}

