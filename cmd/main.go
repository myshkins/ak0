package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
)

func NewServer(
  // logger *Logger
  // config *Config
  // commentStore *commentStore
  // anotherStore *anotherStore
) http.Handler {
  mux := http.NewServeMux()
  // addRoutes()
  var handler http.Handler = mux
  // handler = someMiddleware(handler)
  return handler
}

func run(ctx context.Context, w io.Writer, args []string) error {
  ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
  defer cancel()
  // ...

  return nil
}

func main() {
  ctx := context.Background()
  if err := run(ctx, os.Stdout, os.Args); err != nil {
    fmt.Fprintf(os.Stderr, "%s\n", err)
    os.Exit(1)
  }
}

