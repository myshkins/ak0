package middleware

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func MetricsMiddleware(h http.Handler) http.Handler {

	fn := otelhttp.NewHandler(h, "otel request metrics")
	// fn := func(w http.ResponseWriter, r *http.Request) {
	// 	// call the original http.Handler we're wrapping
	// 	h.ServeHTTP(w, r)

	// 	// record metrics
	// }

	return fn
}
