package logger

import (
  "log/slog"
  "os"
)

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

// todo: implement debug logging, so when in dev info logs go to file, but also log to stdout

