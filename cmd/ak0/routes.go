package main

import (
  "embed"
	"net/http"

	"github.com/myshkins/ak0/internal/handlers"
	"github.com/myshkins/ak0/internal/middleware"
)

func newMiddleware(
	b *middleware.BlockList,
	c *middleware.ClientRateLimiters,
  s *embed.FS,
) func(handler http.Handler) http.Handler {
	m := func(h http.Handler) http.Handler {
		return middleware.LoggingMiddleWare(
			middleware.CheckRateLimit(c,
				middleware.FilterBots(b, h)))
	}
	return m
}

func addRoutes(
	mux *http.ServeMux,
	bl *middleware.BlockList,
	crl *middleware.ClientRateLimiters,
  static *embed.FS,
	// config Config,
	// authProxy *authProxy
) {
	m := newMiddleware(bl, crl, static)
	mux.Handle("/", m(handlers.HandleHome(static)))
	// mux.Handle("/ping", handlers.Ping())
}
