package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/myshkins/ak0/internal/helpers"

	"github.com/x-way/crawlerdetect"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"golang.org/x/time/rate"
)

type BlockedRateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type BlockList struct {
	BlockedRateLimiters map[string]*BlockedRateLimiter // {"client_ipaddr": BlockRateLimiter}
	Mu                  sync.Mutex
}

type contextKey string

const (
  IsBotKey contextKey = "isBot"
)

var maliciousPaths = []string{
  ".aws",
  ".env",
  ".git",
  ".php",
  ".well-known",
  "docker-compose",
  "XDEBUG",
}

var meter = otel.Meter("github.com/myshkins/ak0")

func NewBlockList() *BlockList {
	blockedMap := make(map[string]*BlockedRateLimiter)
	blocklist := BlockList{
		BlockedRateLimiters: blockedMap,
		Mu:                  sync.Mutex{},
	}
	return &blocklist
}

func CleanupBlocklist(ctx context.Context, bl *BlockList) {
  ticker := time.NewTicker(time.Hour * 1)
  defer ticker.Stop()

	for {
    select {
    case <-ticker.C:
      for k, v := range bl.BlockedRateLimiters {
        if time.Since(v.lastSeen) > time.Hour*240 {
          delete(bl.BlockedRateLimiters, k)
        }
      }
    case <-ctx.Done():
      fmt.Println("closing CleanupBlocklist")
      slog.Info("closing CleanupBlocklist")
      return
  }
	}
}


func block(bl *BlockList, ip string) {
	bl.Mu.Lock()
	defer bl.Mu.Unlock()

	if limiter, ok := bl.BlockedRateLimiters[ip]; ok {
		limiter.lastSeen = time.Now()
		return
	}

	r := BlockedRateLimiter{
		limiter:  rate.NewLimiter(0, 0),
		lastSeen: time.Now(),
	}
	bl.BlockedRateLimiters[ip] = &r
	fmt.Println(bl.BlockedRateLimiters)
	return
}

func isBlocked(bl *BlockList, r *http.Request) bool {
	ip := helpers.GetIpAddr(r)
	if _, ok := bl.BlockedRateLimiters[ip]; ok {
		return true
	}
	return false
}

func isMaliciousRequestPath(bl *BlockList, r *http.Request) bool {
  for _, path := range maliciousPaths {
    if strings.Contains(r.URL.Path, path) {
      slog.Info("malicious request detected")
      block(bl, helpers.GetIpAddr(r))
      return true
    }
  }
	return false
}

func FilterBots(bl *BlockList, h http.Handler) http.Handler {
	botCounter, err := meter.Int64Counter(
		"bot.counter",
		metric.WithDescription("number of bot requests detected"),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		msg := fmt.Sprint("failed to create botCounter: ", err)
		slog.Error(msg)
	}
	maliciousBotCounter, err := meter.Int64Counter(
		"maliciousBot.counter",
		metric.WithDescription("number of malicious bot requests blocked"),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		msg := fmt.Sprint("failed to create maliciousBotCounter: ", err)
		slog.Error(msg)
	}

	humanCounter, err := meter.Int64Counter(
		"human.counter",
		metric.WithDescription("number of (probable) human requests served"),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		msg := fmt.Sprint("failed to create humanCounter: ", err)
		slog.Error(msg)
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		ipaddr := helpers.GetIpAddr(r)
		user_agent := r.Header.Get("User-Agent")
		if isBlocked(bl, r) || isMaliciousRequestPath(bl, r) {
			slog.LogAttrs(
				r.Context(),
				slog.LevelInfo,
				"blocked malicious request",
				slog.String(ipaddr, user_agent),
			)
			maliciousBotCounter.Add(r.Context(), 1)
			http.Error(w, "you've been blocked", http.StatusNotFound)
			return
		}

		if crawlerdetect.IsCrawler(user_agent) {
			slog.LogAttrs(
				r.Context(),
				slog.LevelInfo,
				"bot request detected",
				slog.String(ipaddr, user_agent),
			)
			hctx := context.WithValue(r.Context(), IsBotKey, "true")
      r = r.WithContext(hctx)
      botCounter.Add(r.Context(), 1)
		} else {
			slog.LogAttrs(
				r.Context(),
				slog.LevelInfo,
				"request from a possible human served",
				slog.String(ipaddr, user_agent),
			)
      hctx := context.WithValue(r.Context(), IsBotKey, "false")
      r = r.WithContext(hctx)
			humanCounter.Add(r.Context(), 1)
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
