package main

import (
  "net/http"
  "github.com/myshkins/ak0_2/internal/handlers"
)

func addRoutes(
  mux *http.ServeMux,
  // some dependencies here, eg.
  // config Config,
  // authProxy *authProxy
) {
  mux.Handle("/", handlers.HandleHome())
}
