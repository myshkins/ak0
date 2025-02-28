package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
  "os/user"
	"syscall"
	"time"
)

const (
	logPath  = "/ak0/ak0.log"
	logMode  = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	logPerms = 0640
)

func printFileUserInfo() {
  fmt.Printf("\nuser id: %v\ngroup id: %v\n", os.Getuid(), os.Getegid())

  f, err := os.ReadFile("/etc/passwd")
  if err != nil {
    fmt.Printf("error reading /etc/passwd %v\n", err)
  }
  fmt.Println("here's /etc/passwd:")
  fmt.Println(string(f))

  fileInfo, err := os.Stat(logPath)
  if err != nil {
		fmt.Printf("failed to get file info: %w", err)
    return
	}

	// Print basic file information
	fmt.Printf("Mode: %s\n", fileInfo.Mode())

	// Get and print permissions in octal format
	permissions := fileInfo.Mode().Perm()
	fmt.Printf("Permissions: %04o (", permissions)

  if stat, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
		// Get user information
		u, err := user.LookupId(fmt.Sprintf("%d", stat.Uid))
		if err != nil {
			fmt.Printf("Owner UID: %d (unknown user)\n", stat.Uid)
		} else {
			fmt.Printf("Owner: %s (UID: %d)\n", u.Username, stat.Uid)
		}

		// Get group information
		g, err := user.LookupGroupId(fmt.Sprintf("%d", stat.Gid))
		if err != nil {
			fmt.Printf("Group GID: %d (unknown group)\n", stat.Gid)
		} else {
			fmt.Printf("Group: %s (GID: %d)\n", g.Name, stat.Gid)
		}
	} else {
		fmt.Println("Ownership information not available on this platform")
	}	
}

func NewLogger() (*slog.Logger, *os.File) {
  printFileUserInfo()
	out := os.Stdout
	if os.Getenv("AK0_ENV") == "prod" {
		// add retry logic in case of logratate race condition
		for i := 0; i < 5; i++ {
			time.Sleep(time.Duration(i) * 100 * time.Millisecond)
			f, err := os.OpenFile(logPath, logMode, logPerms)
			if err != nil {
        fmt.Fprintf(os.Stderr, "\n %v - failed to open new log file after logrotate. error: %v", time.Now(), err.Error())
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

func ListenForLogrotate(oldfile *os.File, ctx context.Context) {
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
        case syscall.SIGUSR1:
          fmt.Println("logger sigChan received USR1, rotating logs")
          slog.Info("logger sigChan received USR1, rotating logs")
          l, f := NewLogger()
          slog.SetDefault(l)
          err := oldfile.Close()
          if err != nil {
            fmt.Printf("\nak0 Logger: error closing old log file: %v\n", err.Error())
          }
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
