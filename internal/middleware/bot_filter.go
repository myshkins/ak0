package middleware

import (
  "fmt"
	"log/slog"
	"net/http"

	"github.com/myshkins/ak0_2/internal/helpers"

	"github.com/x-way/crawlerdetect"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var meter = otel.Meter("github.com/myshkins/ak0_2")

func FilterBots(h http.Handler) http.Handler {
  botCounter, err := meter.Int64Counter(
		"bot.counter",
		metric.WithDescription("Number of calls for HandleHome"),
		metric.WithUnit("{call}"),
	)
  if err != nil {
    msg := fmt.Sprintf("failed to create botCounter: ", err)
    slog.Error(msg)
  }

  humanCounter, err := meter.Int64Counter(
		"human.counter",
		metric.WithDescription("Number of calls for HandleHome"),
		metric.WithUnit("{call}"),
	)
  if err != nil {
    msg := fmt.Sprintf("failed to create humanCounter: ", err)
    slog.Error(msg)
  }

	fn := func(w http.ResponseWriter, r *http.Request) {
		ipaddr := helpers.GetIpAddr(r)
		user_agent := r.Header.Get("User-Agent")
		if crawlerdetect.IsCrawler(user_agent) {
			slog.LogAttrs(
				r.Context(),
				slog.LevelInfo,
				"bot found",
				slog.String(ipaddr, user_agent),
			)
      botCounter.Add(r.Context(), 1)
		} else {
      humanCounter.Add(r.Context(), 1)
      slog.LogAttrs(
				r.Context(),
				slog.LevelInfo,
				"human request served",
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

