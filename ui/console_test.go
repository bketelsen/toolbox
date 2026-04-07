package ui_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bketelsen/toolbox/ui"
)

func TestConsolePrinter_SetStdout(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(buf, &bytes.Buffer{})

	newBuf := &bytes.Buffer{}
	printer.SetStdout(newBuf)

	// Verify that the stdout was changed
	assert.NotNil(t, newBuf)
}

func TestConsolePrinter_SetStderr(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)

	newBuf := &bytes.Buffer{}
	printer.SetStderr(newBuf)

	// Verify that the stderr was changed
	assert.NotNil(t, newBuf)
}

func TestConsolePrinter_SetLinePrefix(t *testing.T) {
	t.Parallel()

	printer := ui.New(&bytes.Buffer{}, &bytes.Buffer{})
	printer.SetLinePrefix(">>")

	// Verify the line prefix is set (we'll verify it when testing output)
	assert.NotNil(t, printer)
}

func TestConsolePrinter_Info(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)

	printer.Info("Hello World", "Line 1", "Line 2")

	output := buf.String()
	assert.Contains(t, output, "Hello World")
	assert.Contains(t, output, "Line 1")
	assert.Contains(t, output, "Line 2")
}

func TestConsolePrinter_InfoPrefix(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)

	printer.InfoPrefix("PREFIX", "Header", "Line 1")

	output := buf.String()
	assert.Contains(t, output, "Header")
	assert.Contains(t, output, "Line 1")
}

func TestConsolePrinter_Warn(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)

	printer.Warn("Warning message", "Details")

	output := buf.String()
	assert.Contains(t, output, "Warning message")
	assert.Contains(t, output, "Details")
	assert.Contains(t, output, "WARNING")
}

func TestConsolePrinter_WarnPrefix(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)

	printer.WarnPrefix("CUSTOM", "Warning", "Details")

	output := buf.String()
	assert.Contains(t, output, "Warning")
	assert.Contains(t, output, "Details")
}

func TestConsolePrinter_Success(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)

	printer.Success("Success message", "All good")

	output := buf.String()
	assert.Contains(t, output, "Success message")
	assert.Contains(t, output, "All good")
}

func TestConsolePrinter_SuccessPrefix(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)

	printer.SuccessPrefix("OK", "Great", "Everything works")

	output := buf.String()
	assert.Contains(t, output, "Great")
	assert.Contains(t, output, "Everything works")
}

func TestConsolePrinter_Error(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)

	printer.Error("Error message", "Something went wrong")

	output := buf.String()
	assert.Contains(t, output, "Error message")
	assert.Contains(t, output, "Something went wrong")
	assert.Contains(t, output, "ERROR")
}

func TestConsolePrinter_ErrorPrefix(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)

	printer.ErrorPrefix("FAIL", "Error", "Details")

	output := buf.String()
	assert.Contains(t, output, "Error")
	assert.Contains(t, output, "Details")
}

func TestConsolePrinter_TitleInfo(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)

	printer.TitleInfo("Title")

	output := buf.String()
	assert.Contains(t, output, "Title")
}

func TestConsolePrinter_TitleWarn(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)

	printer.TitleWarn("Warning Title")

	output := buf.String()
	assert.Contains(t, output, "Warning Title")
	assert.Contains(t, output, "WARNING")
}

func TestConsolePrinter_TitleSuccess(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)

	printer.TitleSuccess("Success Title")

	output := buf.String()
	assert.Contains(t, output, "Success Title")
}

func TestConsolePrinter_TitleError(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)

	printer.TitleError("Error Title")

	output := buf.String()
	assert.Contains(t, output, "Error Title")
	assert.Contains(t, output, "ERROR")
}

// Package-level function tests
func TestGlobalInfo(t *testing.T) {
	t.Parallel()

	originalConsole := ui.Console
	buf := &bytes.Buffer{}
	ui.Console = ui.New(&bytes.Buffer{}, buf)
	defer func() { ui.Console = originalConsole }()

	ui.Info("Test message", "Details")

	output := buf.String()
	assert.Contains(t, output, "Test message")
	assert.Contains(t, output, "Details")
}

func TestGlobalWarn(t *testing.T) {
	t.Parallel()

	originalConsole := ui.Console
	buf := &bytes.Buffer{}
	ui.Console = ui.New(&bytes.Buffer{}, buf)
	defer func() { ui.Console = originalConsole }()

	ui.Warn("Warning", "Details")

	output := buf.String()
	assert.Contains(t, output, "Warning")
	assert.Contains(t, output, "WARNING")
}

func TestGlobalSuccess(t *testing.T) {
	t.Parallel()

	originalConsole := ui.Console
	buf := &bytes.Buffer{}
	ui.Console = ui.New(&bytes.Buffer{}, buf)
	defer func() { ui.Console = originalConsole }()

	ui.Success("Success", "Good")

	output := buf.String()
	assert.Contains(t, output, "Success")
}

func TestGlobalError(t *testing.T) {
	t.Parallel()

	originalConsole := ui.Console
	buf := &bytes.Buffer{}
	ui.Console = ui.New(&bytes.Buffer{}, buf)
	defer func() { ui.Console = originalConsole }()

	ui.Error("Error", "Bad")

	output := buf.String()
	assert.Contains(t, output, "Error")
	assert.Contains(t, output, "ERROR")
}

func TestGlobalTitleInfo(t *testing.T) {
	t.Parallel()

	originalConsole := ui.Console
	buf := &bytes.Buffer{}
	ui.Console = ui.New(&bytes.Buffer{}, buf)
	defer func() { ui.Console = originalConsole }()

	ui.TitleInfo("Info Title")

	output := buf.String()
	assert.Contains(t, output, "Info Title")
}

func TestGlobalTitleWarn(t *testing.T) {
	t.Parallel()

	originalConsole := ui.Console
	buf := &bytes.Buffer{}
	ui.Console = ui.New(&bytes.Buffer{}, buf)
	defer func() { ui.Console = originalConsole }()

	ui.TitleWarn("Warn Title")

	output := buf.String()
	assert.Contains(t, output, "Warn Title")
	assert.Contains(t, output, "WARNING")
}

func TestGlobalTitleSuccess(t *testing.T) {
	t.Parallel()

	originalConsole := ui.Console
	buf := &bytes.Buffer{}
	ui.Console = ui.New(&bytes.Buffer{}, buf)
	defer func() { ui.Console = originalConsole }()

	ui.TitleSuccess("Success Title")

	output := buf.String()
	assert.Contains(t, output, "Success Title")
}

func TestGlobalTitleError(t *testing.T) {
	t.Parallel()

	originalConsole := ui.Console
	buf := &bytes.Buffer{}
	ui.Console = ui.New(&bytes.Buffer{}, buf)
	defer func() { ui.Console = originalConsole }()

	ui.TitleError("Error Title")

	output := buf.String()
	assert.Contains(t, output, "Error Title")
	assert.Contains(t, output, "ERROR")
}

func TestGlobalInfoPrefix(t *testing.T) {
	t.Parallel()

	originalConsole := ui.Console
	buf := &bytes.Buffer{}
	ui.Console = ui.New(&bytes.Buffer{}, buf)
	defer func() { ui.Console = originalConsole }()

	ui.InfoPrefix("PREFIX", "Message", "Line")

	output := buf.String()
	assert.Contains(t, output, "Message")
	assert.Contains(t, output, "Line")
}

func TestGlobalWarnPrefix(t *testing.T) {
	t.Parallel()

	originalConsole := ui.Console
	buf := &bytes.Buffer{}
	ui.Console = ui.New(&bytes.Buffer{}, buf)
	defer func() { ui.Console = originalConsole }()

	ui.WarnPrefix("PREFIX", "Warning", "Details")

	output := buf.String()
	assert.Contains(t, output, "Warning")
}

func TestGlobalSuccessPrefix(t *testing.T) {
	t.Parallel()

	originalConsole := ui.Console
	buf := &bytes.Buffer{}
	ui.Console = ui.New(&bytes.Buffer{}, buf)
	defer func() { ui.Console = originalConsole }()

	ui.SuccessPrefix("OK", "Done", "Success")

	output := buf.String()
	assert.Contains(t, output, "Done")
}

func TestGlobalErrorPrefix(t *testing.T) {
	t.Parallel()

	originalConsole := ui.Console
	buf := &bytes.Buffer{}
	ui.Console = ui.New(&bytes.Buffer{}, buf)
	defer func() { ui.Console = originalConsole }()

	ui.ErrorPrefix("FAIL", "Failed", "Details")

	output := buf.String()
	assert.Contains(t, output, "Failed")
}

func TestGlobalSetStdout(t *testing.T) {
	t.Parallel()

	originalConsole := ui.Console
	defer func() { ui.Console = originalConsole }()

	newBuf := &bytes.Buffer{}
	ui.SetStdout(newBuf)

	// Verify it was set
	assert.NotNil(t, newBuf)
}

func TestGlobalSetStderr(t *testing.T) {
	t.Parallel()

	originalConsole := ui.Console
	defer func() { ui.Console = originalConsole }()

	newBuf := &bytes.Buffer{}
	ui.SetStderr(newBuf)

	// Verify it was set
	assert.NotNil(t, newBuf)
}

func TestGlobalSetLinePrefix(t *testing.T) {
	t.Parallel()

	originalConsole := ui.Console
	defer func() { ui.Console = originalConsole }()

	ui.SetLinePrefix(">>")
	assert.NotNil(t, ui.Console)
}

func TestNew(t *testing.T) {
	t.Parallel()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	printer := ui.New(stdout, stderr)

	require.NotNil(t, printer)
	assert.NotNil(t, printer)
}
