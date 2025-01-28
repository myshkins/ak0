package middleware

import (
	"log/slog"
	"net/http"

	"github.com/myshkins/ak0_2/internal/helpers"

	"github.com/felixge/httpsnoop"
)

func LoggingMiddleWare(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// this calls the handler as well as capturing metrics
		m := httpsnoop.CaptureMetrics(h, w, r)

		// 200 GET / ipaddr referer mozillafirefoxagent size duration
		slog.LogAttrs(
			r.Context(), // this isn't really used currently, but could implement custom handler to extract context info and log it
			slog.LevelInfo,
			"generic request log",
			slog.Int("status_code", m.Code),
			slog.String("method", r.Method),
			slog.String("uri", r.URL.String()),
			slog.String("ipaddr", helpers.GetIpAddr(r)),
			slog.String("referer", r.Header.Get("Referer")),
			slog.String("user_agent", r.Header.Get("User-Agent")),
			slog.String("ipaddr", helpers.GetIpAddr(r)),
			slog.Int("bytes_written", int(m.Written)),
			slog.Duration("duration", m.Duration),
		)
	}
	return http.HandlerFunc(fn)
}
