package main

import (
  "net/http"

  "github.com/myshkins/ak0/internal/handlers"
	"github.com/myshkins/ak0/internal/middleware"
)

func addRoutes(
  mux *http.ServeMux,
  bl  *middleware.BlockList,
  crl *middleware.ClientRateLimiters,
  // config Config,
  // authProxy *authProxy
) {
  mux.Handle("/",
    middleware.LoggingMiddleWare(
      middleware.CheckRateLimit(crl,
        middleware.FilterBots(bl, handlers.HandleHome()),
        ),
      ),
    )
  // mux.Handle("/ping", handlers.Ping())
}
