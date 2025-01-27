package middleware

import (
	"context"
  "log/slog"
	"net/http"

  "github.com/myshkins/ak0_2/internal/helpers"

  "golang.org/x/time/rate"
)


/*
get ipaddr, check if it's on block or throttle list
check user-agent, add to throttle or block list if conditions met
log bot metrics
*/
var clientLimiters map[string]rate.Limiter

func FilterBots(h http.Handler) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {
    ipaddr := helpers.GetIpAddr(r)
    slog.Info("bot_filter", "ipaddr", ipaddr)
    user_agent := r.Header.Get("User-Agent")
    slog.Info("bot_filter", "user_agent", user_agent)

    ctx := context.Background()
    slog.LogAttrs(
      ctx,
      slog.LevelInfo,
      "",
      slog.String("ipaddr", ipaddr),
      )
}
  return http.HandlerFunc(fn)
}

