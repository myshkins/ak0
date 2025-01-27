package middleware

import (
	"context"
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
    ctx := context.Background()
    slog.LogAttrs(
      ctx,
      slog.LevelInfo,
      "",
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

