package tint

import (
	"io"
	"log/slog"
	"time"
)

var lvl = new(slog.LevelVar)

func Initialize(w io.Writer) {

	lvl.Set(slog.LevelInfo)
	// set global logger with custom options
	slog.SetDefault(slog.New(
		NewHandler(w, &Options{
			Level:      lvl,
			TimeFormat: time.Kitchen,
		}),
	))
}

func SetLevel(level slog.Level) {
	lvl.Set(level)
}
