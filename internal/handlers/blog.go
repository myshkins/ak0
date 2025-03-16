package handlers

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/myshkins/ak0/internal/middleware"
	"go.opentelemetry.io/otel/metric"
)


func HandleBlog(static *embed.FS) http.Handler {
	blogRequestCounter, err := meter.Int64Counter(
		"blogVisit.counter",
		metric.WithDescription("Number of bot/human requests on /"),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		panic(err)
	}

	botBlogRequestCounter, err := meter.Int64Counter(
		"botBlogVisit.counter",
    metric.WithDescription("Number of bot reqeust on /"),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		panic(err)
	}

	humanBlogRequestCounter, err := meter.Int64Counter(
		"humanBlogVisit.counter",
		metric.WithDescription("Number of human requests on /"),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		panic(err)
	}

  staticFiles, err := fs.Sub(*static, "dist")
  if err != nil {
    panic(err)
  }
	fs := http.FileServer(http.FS(staticFiles))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/blog" {
      slog.Info("HandleBlog handling the blog")
      blogRequestCounter.Add(r.Context(), 1)
      if r.Context().Value(middleware.IsBotKey) == "true" {
        botBlogRequestCounter.Add(r.Context(), 1)
      } else {
        humanBlogRequestCounter.Add(r.Context(), 1)
      }
    }
    fs.ServeHTTP(w, r)
  })
}
