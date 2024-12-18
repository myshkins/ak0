package main

import (
  "net/http"
  "log"
  "github.com/myshkins/ak0_2/internal/handlers"
)

func addRoutes(
  mux *http.ServeMux,
  // some dependencies here, eg.
  logger *log.Logger,
  // config Config,
  // authProxy *authProxy
) {
  mux.Handle("/", handlers.HandleHome(logger))
}
