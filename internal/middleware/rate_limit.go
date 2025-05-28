package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/myshkins/ak0/internal/helpers"

	"golang.org/x/time/rate"
)

type ClientRateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}
type ClientRateLimiters struct {
	ClientLimiters map[string]*ClientRateLimiter // {"client_ipaddr": ClientRateLimiter}
	Mu             sync.Mutex
}

func NewClientRateLimiters() *ClientRateLimiters {
	crlMap := make(map[string]*ClientRateLimiter)
	crl := ClientRateLimiters{
		ClientLimiters: crlMap,
		Mu:             sync.Mutex{},
	}
	return &crl
}

func CleanupRateLimiters(ctx context.Context, crl *ClientRateLimiters) {
  ticker := time.NewTicker(time.Minute * 1)
  defer ticker.Stop()

	for {
    select {
    case <-ticker.C:
      for k, v := range crl.ClientLimiters {
        if time.Since(v.lastSeen) > time.Minute*3 {
          delete(crl.ClientLimiters, k)
        }
      }
    case <-ctx.Done():
      slog.Info("closing CleanupRateLimiters")
      return
  }
	}
}

func CheckRateLimit(crl *ClientRateLimiters, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ip := helpers.GetIpAddr(r)

		if !getLimiter(crl, ip).Allow() {
			// todo flag it as bot

			slog.LogAttrs(
				r.Context(),
				slog.LevelInfo,
				"rate_limit",
				slog.String("ipaddr", ip),
			)

			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func getLimiter(crl *ClientRateLimiters, ip string) *rate.Limiter {
	crl.Mu.Lock()
	defer crl.Mu.Unlock()

	if limiter, ok := crl.ClientLimiters[ip]; ok {
		limiter.lastSeen = time.Now()
		return limiter.limiter
	}

	r := ClientRateLimiter{
		limiter:  rate.NewLimiter(20, 100),
		lastSeen: time.Now(),
	}
	crl.ClientLimiters[ip] = &r
	return r.limiter
}
