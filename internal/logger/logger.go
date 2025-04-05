package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	logMode  = os.O_APPEND | os.O_WRONLY
	logPerms = 0640
)


func NewLogger(logPath string) (*os.File) {
	out := os.Stdout
	// if os.Getenv("AK0_ENV") == "prod" {
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

	logger := slog.New(slog.NewJSONHandler(out, nil))
  slog.SetDefault(logger)
	return out
}

func ListenForLogrotate(logPath string, oldfile *os.File, ctx context.Context) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGUSR1)

  done := ctx.Done()

	go func() {
		slog.Info("logger listening for USR1 signal...")
    for {
      select {
      case sig := <-sigChan:
        // either usr1 or sigterm or sigint
        switch sig {
        case syscall.SIGUSR1:  // this isn't used, just restaring container now
          fmt.Println("logger sigChan received USR1, rotating logs")
          slog.Info("logger sigChan received USR1, rotating logs")

          err := oldfile.Close()
          if err != nil {
            fmt.Printf("\nak0 Logger: error closing old log file: %v\n", err.Error())
          } else {
            fmt.Println("successfully closed old file")
          }

          // sleep to allow time for logrotate to finish
          time.Sleep(100 * time.Millisecond)

          f := NewLogger(logPath)
          oldfile = f
        case syscall.SIGTERM, syscall.SIGINT:
          fmt.Println("logger sigChan received sigterm or sigint, shutting down")
          slog.Info("logger sigChan received sigterm or sigint, shutting down")
          signal.Stop(sigChan)
          return
        }
      case <-done:
        fmt.Println("logger context finished, shutting down")
        slog.Info("logger context finished, shutting down")
        signal.Stop(sigChan)
        err := oldfile.Close()
        if err != nil {
          fmt.Printf("\nak0 Logger: error closing old log file: %v\n", err.Error())
        }
        return
      }
    }
  }()
}
