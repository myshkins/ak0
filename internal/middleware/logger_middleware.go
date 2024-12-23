package middleware

import (
  "net/http"

  "github.com/myshkins/ak0_2/internal/logger"
)


func logHTTPReq(uri, method string) {
  logger.Logger.Printf("handledRequest: %s, %s", uri, method)
}

func LoggingMiddleWare(h http.Handler) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {

    // call the original http.Handler we're wrapping
    h.ServeHTTP(w, r)

    // log info
    uri := r.URL.String()
    method := r.Method
    logHTTPReq(uri, method)
  }

  return http.HandlerFunc(fn)
}
