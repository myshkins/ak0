package middleware

import (
	"log/slog"
	"net/http"

	"github.com/myshkins/ak0_2/internal/helpers"

	"github.com/x-way/crawlerdetect"
)

func FilterBots(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ipaddr := helpers.GetIpAddr(r)
		user_agent := r.Header.Get("User-Agent")
		if crawlerdetect.IsCrawler(user_agent) {
			slog.Info("bot_filter", ipaddr, "is a bot")
			slog.LogAttrs(
				r.Context(),
				slog.LevelInfo,
				"bot found",
				slog.String(ipaddr, user_agent),
			)
		}
		/*
		   TODO: implement additonal bot filtering
		   check for "POST" request
		   header order
		   ip range, look for unusual location
		   check for high rate of requests
		   js cookie challenge?
		*/
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
