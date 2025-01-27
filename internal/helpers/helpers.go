package helpers

import (
  "net/http"
)

func GetIpAddr(r *http.Request) string {
  // default to using custom set header
  ip := r.Header.Get("AK-First-External-IP")
  if ip != "" {
    return ip
  }

  // if above fails for some reason, use RemnoteAddr
  // note, that will just be the ip of the nginx reverse proxy
  ip = r.RemoteAddr
  return ip
}
