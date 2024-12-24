package middleware

import (
	"context"
  "log/slog"
	"net/http"

	"github.com/myshkins/ak0_2/internal/logger"
)


func LoggingMiddleWare(h http.Handler) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {

    // call the original http.Handler we're wrapping
    h.ServeHTTP(w, r)

    // log info
    // 200 GET / referer mozillafirefoxagent size duration
    ctx := context.Background()
    logger.Logger.LogAttrs(
      ctx,
      slog.LevelInfo,
      "",
      slog.Int("status_code", 200),
      slog.String("method", r.Method),
      slog.String("uri", r.URL.String()),
      slog.String("referer", r.Header.Get("Referer")),
      slog.String("user_agent", r.Header.Get("User-Agent")),
      slog.String("ipaddr", getIpAddr()), 
      )
}
  return http.HandlerFunc(fn)
}

func getIpAddr() string {
  ip := ""
  return ip
}
