package middleware

import (
	"context"
  "log/slog"
	"net/http"

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
      slog.String("ipaddrg", getIpAddr(r)),
      slog.String("referer", r.Header.Get("Referer")),
      slog.String("user_agent", r.Header.Get("User-Agent")),
      slog.String("ipaddr", getIpAddr(r)),
      slog.Int("bytes_written", int(m.Written)),
      slog.Duration("duration", m.Duration),
      )
}
  return http.HandlerFunc(fn)
}

func getIpAddr(r *http.Request) string {
  // default to using custom set header
  ip := r.Header.Get("AK-First-External-IP")
  if ip != "" {
    return ip
  }

  // if above fails for some reason, use RemnoteAddr
  // note, that will just be the ip of the nginx reverse proxy
  ip = r.RemoteAddr
  return ip
}
