package handlers

import (
	"log/slog"
	"net/http"
)

func HandleHome(logger *slog.Logger) http.Handler {
  // do necessary prep work here
  logger.Info("HandleHome was called")
  fs := http.FileServer(http.Dir("/home/myshkins/projects/ak0_2/web/dist"))
  fsWrapper := func() http.Handler {
    logger.Info("fsWrapper was called")
    return fs
  }
  return fsWrapper()
}

