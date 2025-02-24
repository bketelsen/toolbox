package tint

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

var lvl = new(slog.LevelVar)

func Initialize(w io.Writer, debug bool) {

	exe, _ := os.Executable()
	debugEnv := fmt.Sprintf("%s_DEBUG", strings.ToUpper(exe))

	_, ok := os.LookupEnv(debugEnv)
	if ok {
		debug = true
	}
	if debug {
		lvl.Set(slog.LevelDebug)
	} else {
		lvl.Set(slog.LevelInfo)

	}
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
