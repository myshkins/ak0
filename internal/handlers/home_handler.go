package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func HandleHome(logger *log.Logger) http.Handler {
  // do necessary prep work here
  log.Println("HandleHome was called")
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
      fmt.Fprintf(w, "hewo myshy! you requested: %s\n", r.URL.Path)
      logger.Println("handling home")
		},
	)
}

