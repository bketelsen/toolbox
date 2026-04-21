package toolbox

import (
	"log/slog"
	"os"
)

// setupLogger initializes a.Logger based on the current flag state.
// It opens a.LogFile for appending if set, otherwise falls back to os.Stderr.
// The handler type is determined by the JSONOutput flag.
func (a *App) setupLogger() error {
	level := slog.LevelInfo
	if Verbose {
		level = slog.LevelDebug
	}
	opts := &slog.HandlerOptions{Level: level}

	var w *os.File
	if a.LogFile != "" {
		f, err := os.OpenFile(a.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		a.logFile = f
		w = f
	} else {
		w = os.Stderr
	}

	var handler slog.Handler
	if JSONOutput {
		handler = slog.NewJSONHandler(w, opts)
	} else {
		handler = slog.NewTextHandler(w, opts)
	}

	a.Logger = slog.New(handler)
	return nil
}

// Close closes the log file if one was opened. Safe to call multiple times.
func (a *App) Close() error {
	if a.logFile != nil {
		return a.logFile.Close()
	}
	return nil
}
