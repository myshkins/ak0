package middleware

import "net/http"

func MetricsMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		// call the original http.Handler we're wrapping
		h.ServeHTTP(w, r)

		// record metrics

	}

	return http.HandlerFunc(fn)
}
