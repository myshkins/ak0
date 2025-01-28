package main

import (
  "net/http"

  "github.com/myshkins/ak0_2/internal/handlers"
	"github.com/myshkins/ak0_2/internal/middleware"
)

func addRoutes(
  mux *http.ServeMux,
  // some dependencies here, eg.
  // config Config,
  // authProxy *authProxy
) {
  mux.Handle("/",
    middleware.FilterBots(
      middleware.CheckRateLimit(
        middleware.LoggingMiddleWare(handlers.HandleHome()),
        ),
      ),
    )
}
