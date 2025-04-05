package helpers

import (
  "fmt"
	"net/http"
  "path/filepath"
  "runtime"
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

// turns a relative path into an absolute one
func MakeRelPathAbs(relativePath string) (string, error) {
  _, filename, _, ok := runtime.Caller(1)
  if !ok {
    return "", fmt.Errorf("error getting source file path")
  }

  sourceDir := filepath.Dir(filename)
  return filepath.Join(sourceDir, relativePath), nil
}
