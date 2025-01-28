package logger

import (
	"log"
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	out := os.Stdout
	if os.Getenv("AK0_2_ENV") == "prod" {
		f, err := os.OpenFile("/ak0_2_log", os.O_RDWR, os.ModeAppend)
		if err != nil {
			log.Fatal(err)
		}
		out = f
	}
	logger := slog.New(slog.NewJSONHandler(out, nil))
	return logger
}
