package handlers

import (
  "context"
	"log/slog"
	"net/http"
	"os"

  "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var meter = otel.Meter("github.com/myshkins/ak0")


func HandleHome() http.Handler {
	// do necessary prep work here
	slog.Info("HandleHome was called")

	apiCounter, err := meter.Int64Counter(
		"test.counter",
		metric.WithDescription("Number of calls for HandleHome"),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		panic(err)
	}
  ctx := context.Background()
  c := context.WithValue(ctx, "ak02contextest", "hewo")

  fp := "/home/myshkins/projects/ak0/web/dist"
	if os.Getenv("AK0_ENV") == "prod" {
		fp = "/lib/node_modules/ak02/dist"
	}
	fs := http.FileServer(http.Dir(fp))

	fsWrapper := func() http.Handler {
		slog.Info("fsWrapper was called")
    apiCounter.Add(c, 1)
		return fs
	}
	return fsWrapper()
}
