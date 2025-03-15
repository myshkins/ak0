package main

import (
	"net/http"

	"github.com/myshkins/ak0/internal/handlers"
	"github.com/myshkins/ak0/internal/middleware"
)

func newMiddleware(
	b *middleware.BlockList,
	c *middleware.ClientRateLimiters,
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
	// config Config,
	// authProxy *authProxy
) {
	m := newMiddleware(bl, crl)
	mux.Handle("/", m(handlers.HandleHome()))
	// mux.Handle("/ping", handlers.Ping())
}
