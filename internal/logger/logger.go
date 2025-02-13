package logger

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	logPath  = "/ak0.log"
	logMode  = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	logPerms = 0640
)

func NewLogger() (*slog.Logger, *os.File) {
	// out := os.Stdout
  out, err := os.OpenFile("/home/myshkins/projects/ak0/dev.log", logMode, logPerms)
  if err != nil {fmt.Println(err)}
	if os.Getenv("AK0_ENV") == "prod" {
		// add retry logic in case of logratate race condition
		for i := 0; i < 5; i++ {
			time.Sleep(100 * time.Millisecond)
			f, err := os.OpenFile(logPath, logMode, logPerms)
			if err != nil {
				fmt.Fprintf(os.Stderr, "\n %v - failed to open new log file after logrotate", time.Now())
				continue
			}
			out = f
			break
		}

	}
	logger := slog.New(slog.NewJSONHandler(out, nil))
  slog.SetDefault(logger)
	return logger, out
}

func ListenForLogrotate(oldfile *os.File) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGUSR1)

	go func() {
		slog.Info("logger listening for USR1 signal...")
		for sig := range sigChan {
			switch sig {
			case syscall.SIGUSR1:
				fmt.Println("logger sigChan received USR1, rotating logs")
				l, f := NewLogger()
				slog.SetDefault(l)
				oldfile.Close()
				oldfile = f
			case syscall.SIGTERM:
				fmt.Println("logger sigChan received sigterm, shutting down")
				return
			}
		}
	}()
}
