package middleware

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/myshkins/ak0_2/internal/helpers"

	"golang.org/x/time/rate"
)

/*
get ipaddr, check if it's on block or throttle list
check user-agent, add to throttle or block list if conditions met
log bot metrics
*/
var clientLimiters = make(map[string]*rate.Limiter) // {"client_ipaddr": limiter}
var mu sync.Mutex

func CheckRateLimit(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ip := helpers.GetIpAddr(r)
		slog.Info("rate_limit", "ipaddr", ip)

		if !getLimiter(ip).Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		h.ServeHTTP(w, r)

		// slog.LogAttrs(
		// 	r.Context(),
		// 	slog.LevelInfo,
		// 	"",
		// 	slog.String("ipaddr", ip),
		// )
	}
	return http.HandlerFunc(fn)
}

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if limiter, ok := clientLimiters[ip]; ok {
		return limiter
	}

	clientLimiters[ip] = rate.NewLimiter(20, 100)
	return clientLimiters[ip]
}
