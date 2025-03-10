// https://github.com/coder/coder/blob/main/LICENSE
// Extracted and modified from github.com/coder/coder
package cobra

import (
	"time"

	"github.com/muesli/termenv"

	"github.com/charmbracelet/lipgloss"
)

// DefaultStyles compose visual elements of the UI.
var DefaultStyles Styles

var (
	BoldStyle = lipgloss.NewStyle().Bold(true)
)

type Styles struct {
	Title,
	Code,
	DateTimeStamp,
	Error,
	Field,
	Hyperlink,
	Keyword,
	Placeholder,
	Prompt,
	FocusedPrompt,
	Fuchsia,
	Warn,
	Success,
	Info,
	Debug,
	Wrap lipgloss.Style
}

var (
	color termenv.Profile
)

var (
	// ANSI color codes
	red           = lipgloss.Color("1")
	green         = lipgloss.Color("2")
	yellow        = lipgloss.Color("3")
	magenta       = lipgloss.Color("5")
	brightBlue    = lipgloss.Color("12")
	brightMagenta = lipgloss.Color("13")
	cyan          = lipgloss.Color("36")
)

func isTerm() bool {
	return color != termenv.Ascii
}

func DisableColor() {
	color = termenv.Ascii
}
func EnableColor() {
	color = termenv.ANSI
}

// Bold returns a formatter that renders text in bold
// if the terminal supports it.
func Bold(s string) string {
	if !isTerm() {
		return s
	}
	return BoldStyle.Render(s)
}

// Bold returns a formatter that renders text in bold
// if the terminal supports it.
func Title(s string) string {
	if !isTerm() {
		return s
	}
	return DefaultStyles.Info.Render(s)
}

// Timestamp formats a timestamp for display.
func Timestamp(t time.Time) string {
	if !isTerm() {
		return t.Format(time.Stamp)
	}
	return DefaultStyles.DateTimeStamp.Render(t.Format(time.Stamp))
}

// Keyword formats a keyword for display.
func Keyword(s string) string {
	if !isTerm() {
		return s
	}
	return DefaultStyles.Keyword.Render(s)
}

// Placeholder formats a placeholder for display.
func Placeholder(s string) string {
	if !isTerm() {
		return s
	}
	return DefaultStyles.Placeholder.Render(s)
}

// Wrap prevents the text from overflowing the terminal.
func Wrap(s string) string {
	if !isTerm() {
		return s
	}
	return DefaultStyles.Wrap.Render(s)
}

// Code formats code for display.
func Code(s string) string {
	if !isTerm() {
		return s
	}
	return DefaultStyles.Code.Render(s)
}

// Field formats a field for display.
func Field(s string) string {
	if !isTerm() {
		return s
	}
	return DefaultStyles.Field.Render(s)
}

// KeyValuePair formats a kvp for display.
func KeyValuePair(key, value string) string {
	if !isTerm() {
		return key + ":" + value
	}
	k := Field(key)
	v := Keyword(value)
	return k + ":" + v
}

func init() {
	// We do not adapt the color based on whether the terminal is light or dark.
	// Doing so would require a round-trip between the program and the terminal
	// due to the OSC query and response.
	DefaultStyles = Styles{
		Title: lipgloss.NewStyle().
			Foreground(brightBlue),
		Code: lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Foreground(lipgloss.Color("#ED567A")).
			Background(lipgloss.Color("#2C2C2C")),
		DateTimeStamp: lipgloss.NewStyle().
			Foreground(brightBlue),

		Error: lipgloss.NewStyle().
			Foreground(red),

		Field: lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#2B2A2A")),

		Fuchsia: lipgloss.NewStyle().
			Foreground(brightMagenta),

		Hyperlink: lipgloss.NewStyle().
			Foreground(magenta).
			Underline(true),

		Keyword: lipgloss.NewStyle().
			Foreground(green),

		Placeholder: lipgloss.NewStyle().
			Foreground(magenta),

		Warn: lipgloss.NewStyle().
			Foreground(yellow),

		Success: lipgloss.NewStyle().
			Foreground(green),

		Info: lipgloss.NewStyle().
			Foreground(brightBlue),
		Debug: lipgloss.NewStyle().
			Foreground(cyan),

		Wrap: lipgloss.NewStyle().
			Width(80),
	}
}
