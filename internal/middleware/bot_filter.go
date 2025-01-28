package middleware

import (
	"log/slog"
	"net/http"

	"github.com/myshkins/ak0_2/internal/helpers"
)

func FilterBots(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ipaddr := helpers.GetIpAddr(r)
		slog.Info("bot_filter", "ipaddr", ipaddr)
		user_agent := r.Header.Get("User-Agent")
		slog.Info("bot_filter", "user_agent", user_agent)

		// slog.LogAttrs(
		// 	r.Context(),
		// 	slog.LevelInfo,
		// 	"",
		// 	slog.String("ipaddr", ipaddr),
		// )
    h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
