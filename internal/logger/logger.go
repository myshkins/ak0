package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

const (
	logMode  = os.O_APPEND | os.O_WRONLY
	logPerms = 0640
)


func NewLogger(logPath string) (*os.File) {
	out := os.Stdout
  if logPath != "stdout" {
    // retry logic in case of logratate race condition
    for i := 0; i < 5; i++ {
      time.Sleep(time.Duration(i) * 100 * time.Millisecond)
      f, err := os.OpenFile(logPath, logMode, logPerms)
      if err == nil {
        fmt.Fprintf(os.Stdout, "\nopened new log file without error\n")
      }
      if err != nil && i == 4 {
        fmt.Fprintf(os.Stdout, "\n %v - failed to open new log file after logrotate. error: %v", time.Now(), err.Error())
        fmt.Fprintf(os.Stdout, "\n this was the last try. using stdout")
        out = os.Stdout
      }
      if err != nil {
        fmt.Fprintf(os.Stdout, "\n %v - failed to open new log file after logrotate. error: %v", time.Now(), err.Error())
        continue
      }
      out = f
      break
    }
  }

	logger := slog.New(slog.NewJSONHandler(out, nil))
  slog.SetDefault(logger)
	return out
}

func ListenForLogrotate(logPath string, oldfile *os.File, ctx context.Context) {
  slog.Info("logger listening for signal...")

  <- ctx.Done()
  time.Sleep(time.Second * 1)
  slog.Info("closing ListenForLogrotate")
  err := oldfile.Close()
  if err != nil {
    fmt.Printf("\nak0 Logger: error closing old log file: %v\n", err.Error())
  }
  return
}
