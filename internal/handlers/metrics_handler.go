package handlers

import (
	"net/http"
)

func HandleMetrics() http.Handler {
	// do necessary prep work here
  mhWrapper := func() http.Handler {
    var mh http.Handler
    return mh
	}
	return mhWrapper()
}
