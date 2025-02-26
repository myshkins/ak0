package handlers

import (
	"net/http"
	"os"

  "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

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

  fp := "/home/myshkins/projects/ak0/web/dist"
	if os.Getenv("AK0_ENV") == "prod" {
		fp = "/lib/node_modules/ak02/dist"
	}
	fs := http.FileServer(http.Dir(fp))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/" {
      homeRequestCounter.Add(r.Context(), 1)
      if r.Context().Value("isBot") == "true" {
        botHomeRequestCounter.Add(r.Context(), 1)
      } else {
        humanHomeRequestCounter.Add(r.Context(), 1)
      }
    }
    fs.ServeHTTP(w, r)
  })
}
