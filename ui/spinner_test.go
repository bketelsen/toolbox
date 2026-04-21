package ui_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"charm.land/lipgloss/v2"
	"github.com/bketelsen/toolbox/ui"
)

func TestNewSpinner(t *testing.T) {
	t.Parallel()

	s := ui.NewSpinner()
	require.NotNil(t, s)
	// Spinner is created successfully
	assert.NotNil(t, s)
}

func TestSpinnerType(t *testing.T) {
	t.Parallel()

	s := ui.NewSpinner()
	result := s.Type(ui.Dots)
	assert.Same(t, s, result, "Type should return same Spinner for chaining")
}

func TestSpinnerTitle(t *testing.T) {
	t.Parallel()

	s := ui.NewSpinner()
	result := s.Title("Custom Title")
	assert.Same(t, s, result, "Title should return same Spinner for chaining")
}

func TestSpinnerOutput(t *testing.T) {
	t.Parallel()

	s := ui.NewSpinner()
	buf := &bytes.Buffer{}
	result := s.Output(buf)
	assert.Same(t, s, result, "Output should return same Spinner for chaining")
}

func TestSpinnerAction(t *testing.T) {
	t.Parallel()

	action := func() {
		// Action implementation
	}

	s := ui.NewSpinner()
	result := s.Action(action)
	assert.Same(t, s, result, "Action should return same Spinner for chaining")
	assert.NotNil(t, s)
}

func TestSpinnerActionWithErr(t *testing.T) {
	t.Parallel()

	action := func(ctx context.Context) error {
		return nil
	}

	s := ui.NewSpinner()
	result := s.ActionWithErr(action)
	assert.Same(t, s, result, "ActionWithErr should return same Spinner for chaining")
	assert.NotNil(t, s)
}

func TestSpinnerContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := ui.NewSpinner()
	result := s.Context(ctx)
	assert.Same(t, s, result, "Context should return same Spinner for chaining")
}

func TestSpinnerStyle(t *testing.T) {
	t.Parallel()

	style := lipgloss.NewStyle().Bold(true)
	s := ui.NewSpinner()
	result := s.Style(style)
	assert.Same(t, s, result, "Style should return same Spinner for chaining")
}

func TestSpinnerTitleStyle(t *testing.T) {
	t.Parallel()

	style := lipgloss.NewStyle().Bold(true)
	s := ui.NewSpinner()
	result := s.TitleStyle(style)
	assert.Same(t, s, result, "TitleStyle should return same Spinner for chaining")
}

func TestSpinnerAccessible(t *testing.T) {
	t.Parallel()

	s := ui.NewSpinner()
	result := s.Accessible(true)
	assert.Same(t, s, result, "Accessible should return same Spinner for chaining")
}

func TestSpinnerInit(t *testing.T) {
	t.Parallel()

	s := ui.NewSpinner()
	cmd := s.Init()
	assert.NotNil(t, cmd, "Init should return a Cmd")
}

func TestSpinnerInitWithAction(t *testing.T) {
	t.Parallel()

	s := ui.NewSpinner().
		ActionWithErr(func(ctx context.Context) error {
			return nil
		})

	cmd := s.Init()
	assert.NotNil(t, cmd, "Init should return a Cmd")
}

func TestSpinnerView(t *testing.T) {
	t.Parallel()

	s := ui.NewSpinner()
	view := s.View()
	assert.NotNil(t, view, "View should return a tea.View")
}

func TestSpinnerViewWithTitle(t *testing.T) {
	t.Parallel()

	s := ui.NewSpinner().Title("Test Title")
	view := s.View()
	assert.NotNil(t, view, "View should return a tea.View")
}

func TestSpinnerBuilderChaining(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	ctx := context.Background()

	s := ui.NewSpinner().
		Title("Processing...").
		Type(ui.Dots).
		Output(buf).
		Context(ctx).
		Accessible(true).
		TitleStyle(lipgloss.NewStyle().Bold(true))

	assert.NotNil(t, s)
}

func TestSpinnerRunWithCancelledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	s := ui.NewSpinner().Context(ctx)
	err := s.Run()
	assert.Error(t, err, "Run should return error when context is cancelled")
}

func TestSpinnerTypeConstants(t *testing.T) {
	t.Parallel()

	// Verify that type constants exist and are accessible
	types := []ui.Type{
		ui.Line,
		ui.Dots,
		ui.MiniDot,
		ui.Jump,
		ui.Points,
		ui.Pulse,
		ui.Globe,
		ui.Moon,
		ui.Monkey,
		ui.Meter,
		ui.Hamburger,
		ui.Ellipsis,
	}

	assert.Equal(t, 12, len(types), "Should have 12 spinner types defined")
}
