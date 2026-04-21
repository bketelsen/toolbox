package toolbox

import (
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestSetupLogger_DefaultsToStderr(t *testing.T) {
	app := &App{}
	if err := app.setupLogger(); err != nil {
		t.Fatalf("setupLogger() error = %v", err)
	}
	if app.Logger == nil {
		t.Fatal("Logger should not be nil")
	}
	if _, ok := app.Logger.Handler().(*slog.TextHandler); !ok {
		t.Errorf("handler type = %T, want *slog.TextHandler", app.Logger.Handler())
	}
}

func TestSetupLogger_JSONHandler(t *testing.T) {
	JSONOutput = true
	defer func() { JSONOutput = false }()

	app := &App{}
	if err := app.setupLogger(); err != nil {
		t.Fatalf("setupLogger() error = %v", err)
	}
	if _, ok := app.Logger.Handler().(*slog.JSONHandler); !ok {
		t.Errorf("handler type = %T, want *slog.JSONHandler", app.Logger.Handler())
	}
}

func TestSetupLogger_TextHandler(t *testing.T) {
	JSONOutput = false

	app := &App{}
	if err := app.setupLogger(); err != nil {
		t.Fatalf("setupLogger() error = %v", err)
	}
	if _, ok := app.Logger.Handler().(*slog.TextHandler); !ok {
		t.Errorf("handler type = %T, want *slog.TextHandler", app.Logger.Handler())
	}
}

func TestSetupLogger_WritesToFile(t *testing.T) {
	logPath := t.TempDir() + "/test.log"
	app := &App{LogFile: logPath}
	defer func() { _ = app.Close() }()

	if err := app.setupLogger(); err != nil {
		t.Fatalf("setupLogger() error = %v", err)
	}

	app.Logger.Info("hello from test")
	_ = app.Close()

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !strings.Contains(string(data), "hello from test") {
		t.Errorf("log file content %q does not contain expected message", string(data))
	}
}

func TestSetupLogger_AppendsToExistingFile(t *testing.T) {
	logPath := t.TempDir() + "/existing.log"
	existing := "pre-existing content\n"
	if err := os.WriteFile(logPath, []byte(existing), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	app := &App{LogFile: logPath}
	if err := app.setupLogger(); err != nil {
		t.Fatalf("setupLogger() error = %v", err)
	}
	app.Logger.Info("appended line")
	_ = app.Close()

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	content := string(data)
	if !strings.HasPrefix(content, existing) {
		t.Errorf("original content not preserved; got %q", content)
	}
	if !strings.Contains(content, "appended line") {
		t.Errorf("appended line not found in %q", content)
	}
}

func TestClose_NoLogFile(t *testing.T) {
	app := &App{}
	if err := app.Close(); err != nil {
		t.Errorf("Close() with no log file = %v, want nil", err)
	}
}

func TestClose_ClosesFile(t *testing.T) {
	logPath := t.TempDir() + "/close.log"
	app := &App{LogFile: logPath}

	if err := app.setupLogger(); err != nil {
		t.Fatalf("setupLogger() error = %v", err)
	}
	if err := app.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
	// After close, writing should fail
	app.Logger.Info("after close")
	// The log call itself won't error (slog swallows write errors),
	// but the file handle is closed. Verify by trying to write directly.
	_, err := app.logFile.WriteString("direct write after close\n")
	if err == nil {
		t.Error("expected error writing to closed file, got nil")
	}
}
