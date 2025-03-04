package handlers

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/myshkins/ak0/internal/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

//go:embed dist/*
var dist embed.FS

var meter = otel.Meter("github.com/myshkins/ak0")

func HandleHome() http.Handler {
	homeRequestCounter, err := meter.Int64Counter(
		"homeVisit.counter",
		metric.WithDescription("Number of bot/human requests on /"),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		panic(err)
	}

	botHomeRequestCounter, err := meter.Int64Counter(
		"botHomeVisit.counter",
    metric.WithDescription("Number of bot reqeust on /"),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		panic(err)
	}

	humanHomeRequestCounter, err := meter.Int64Counter(
		"humanHomeVisit.counter",
		metric.WithDescription("Number of human requests on /"),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		panic(err)
	}

  staticFiles, err := fs.Sub(dist, "dist")
  if err != nil {
    panic(err)
  }
	fs := http.FileServer(http.FS(staticFiles))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/" {
      homeRequestCounter.Add(r.Context(), 1)
      if r.Context().Value(middleware.IsBotKey) == "true" {
        botHomeRequestCounter.Add(r.Context(), 1)
      } else {
        humanHomeRequestCounter.Add(r.Context(), 1)
      }
    }
    fs.ServeHTTP(w, r)
  })
}
