package ui_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bketelsen/toolbox/ui"
)

func TestWarn(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)
	printer.Warn("test header", "line 1", "line 2")

	result := buf.String()
	require.NotEmpty(t, result)
	assert.Contains(t, result, "WARNING")
	assert.Contains(t, result, "test header")
	assert.Contains(t, result, "line 1")
	assert.Contains(t, result, "line 2")
}

func TestWarnPrefix(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)
	printer.WarnPrefix("prefix", "test header", "line 1")

	result := buf.String()
	require.NotEmpty(t, result)
	assert.Contains(t, result, "prefix")
	assert.Contains(t, result, "test header")
}

func TestInfo(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)
	printer.Info("test header", "line 1", "line 2")

	result := buf.String()
	require.NotEmpty(t, result)
	assert.Contains(t, result, "test header")
	assert.Contains(t, result, "line 1")
	assert.Contains(t, result, "line 2")
}

func TestInfoPrefix(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)
	printer.InfoPrefix("prefix", "test header", "line 1")

	result := buf.String()
	require.NotEmpty(t, result)
	assert.Contains(t, result, "prefix")
	assert.Contains(t, result, "test header")
}

func TestSuccess(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)
	printer.Success("test header", "line 1")

	result := buf.String()
	require.NotEmpty(t, result)
	assert.Contains(t, result, "test header")
	assert.Contains(t, result, "line 1")
}

func TestSuccessPrefix(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)
	printer.SuccessPrefix("prefix", "test header", "line 1")

	result := buf.String()
	require.NotEmpty(t, result)
	assert.Contains(t, result, "prefix")
	assert.Contains(t, result, "test header")
}

func TestError(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)
	printer.Error("test header", "line 1")

	result := buf.String()
	require.NotEmpty(t, result)
	assert.Contains(t, result, "ERROR")
	assert.Contains(t, result, "test header")
	assert.Contains(t, result, "line 1")
}

func TestErrorPrefix(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)
	printer.ErrorPrefix("prefix", "test header", "line 1")

	result := buf.String()
	require.NotEmpty(t, result)
	assert.Contains(t, result, "prefix")
	assert.Contains(t, result, "test header")
}

func TestMessageWithoutLines(t *testing.T) {
	t.Parallel()

	t.Run("Warn", func(t *testing.T) {
		t.Parallel()
		buf := &bytes.Buffer{}
		printer := ui.New(&bytes.Buffer{}, buf)
		printer.Warn("header only")

		result := buf.String()
		require.NotEmpty(t, result)
		assert.Contains(t, result, "header only")
	})

	t.Run("Info", func(t *testing.T) {
		t.Parallel()
		buf := &bytes.Buffer{}
		printer := ui.New(&bytes.Buffer{}, buf)
		printer.Info("header only")

		result := buf.String()
		require.NotEmpty(t, result)
		assert.Contains(t, result, "header only")
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		buf := &bytes.Buffer{}
		printer := ui.New(&bytes.Buffer{}, buf)
		printer.Success("header only")

		result := buf.String()
		require.NotEmpty(t, result)
		assert.Contains(t, result, "header only")
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()
		buf := &bytes.Buffer{}
		printer := ui.New(&bytes.Buffer{}, buf)
		printer.Error("header only")

		result := buf.String()
		require.NotEmpty(t, result)
		assert.Contains(t, result, "header only")
	})
}

func TestMessageWithEmptyLines(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	printer := ui.New(&bytes.Buffer{}, buf)
	printer.Info("header", "", "line 1", "", "line 2")

	result := buf.String()
	require.NotEmpty(t, result)
	assert.Contains(t, result, "header")
	assert.Contains(t, result, "line 1")
	assert.Contains(t, result, "line 2")

	// Count newlines to verify empty lines are skipped
	lines := strings.Split(result, "\n")
	nonEmptyLines := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines++
		}
	}
	assert.GreaterOrEqual(t, nonEmptyLines, 3, "Should have at least header and 2 content lines")
}

func TestLinePrefixes(t *testing.T) {
	t.Parallel()

	t.Run("LinePrefixDefault", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "|", ui.LinePrefixDefault)
	})

	t.Run("LinePrefixBullet", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "•", ui.LinePrefixBullet)
	})

	t.Run("LinePrefixCheck", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "✓", ui.LinePrefixCheck)
	})

	t.Run("LinePrefixCross", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "✗", ui.LinePrefixCross)
	})

	t.Run("LinePrefixWarning", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "⚠", ui.LinePrefixWarning)
	})

	t.Run("LinePrefixError", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "✖", ui.LinePrefixError)
	})

	t.Run("LinePrefixSuccess", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "✔", ui.LinePrefixSuccess)
	})

	t.Run("LinePrefixQuestion", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "?", ui.LinePrefixQuestion)
	})

	t.Run("LinePrefixDoubleAngle", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "»", ui.LinePrefixDoubleAngle)
	})
}
