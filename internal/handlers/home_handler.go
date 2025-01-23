package handlers

import (
	"log/slog"
	"net/http"
	"os"
)

func HandleHome(logger *slog.Logger) http.Handler {
  // do necessary prep work here
  logger.Info("HandleHome was called")

  fp := "/home/myshkins/projects/ak0_2/web/dist"
  if os.Getenv("AK0_2_ENV") == "prod" {
    fp = "/lib/node_modules/ak02/dist"
  }
  fs := http.FileServer(http.Dir(fp))

  fsWrapper := func() http.Handler {
    logger.Info("fsWrapper was called")
    return fs
  }
  return fsWrapper()
}

