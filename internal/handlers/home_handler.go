package handlers

import (
	"log"
	"net/http"
)

func HandleHome(logger *log.Logger) http.Handler {
  // do necessary prep work here
  log.Println("HandleHome was called")
  fs := http.FileServer(http.Dir("/home/myshkins/projects/ak0_2/web/dist"))
  fsWrapper := func() http.Handler {
    log.Println("fsWrapper was called")
    return fs
  }
  return fsWrapper()
}

