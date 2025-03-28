// https://github.com/coder/coder/blob/main/LICENSE
// Extracted and modified from github.com/coder/coder
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// cliMessage provides a human-readable message for CLI errors and messages.
type cliMessage struct {
	Style  lipgloss.Style
	Header string
	Prefix string
	Lines  []string
}

// String formats the CLI message for consumption by a human.
func (m cliMessage) String() string {
	var str strings.Builder

	if m.Prefix != "" {
		_, _ = str.WriteString(Bold(m.Prefix))
	}

	str.WriteString(m.Style.Render(m.Header))
	_, _ = str.WriteString("\r\n")
	for _, line := range m.Lines {
		_, _ = fmt.Fprintf(&str, "  %s %s\r\n", m.Style.Render("*"), line)
	}
	return str.String()
}

// Warn writes a log to the writer provided.
func Warn(header string, lines ...string) string {
	return cliMessage{
		Style:  DefaultStyles.Warn,
		Prefix: "WARNING: ",
		Header: header,
		Lines:  lines,
	}.String()
}

func WarnPrefix(prefix, header string, lines ...string) string {
	return cliMessage{
		Style:  DefaultStyles.Warn,
		Prefix: Yellow(prefix) + ": ",
		Header: header,
		Lines:  lines,
	}.String()
}

// Info writes a log to the writer provided.
func Info(header string, lines ...string) string {
	return cliMessage{
		Header: header,
		Lines:  lines,
	}.String()
}

// InfoPrefix writes a log with a prefix to the writer provided.
func InfoPrefix(prefix, header string, lines ...string) string {
	return cliMessage{
		Header: header,
		Prefix: prefix + ": ",
		Lines:  lines,
	}.String()
}

// SuccessPrefix writes a log with a prefix to the writer provided.
func SuccessPrefix(prefix, header string, lines ...string) string {
	return cliMessage{
		Header: header,
		Prefix: Green(prefix) + ": ",
		Lines:  lines,
	}.String()
}

// Error writes a log to the writer provided.
func Error(header string, lines ...string) string {
	return cliMessage{
		Style:  DefaultStyles.Error,
		Prefix: "ERROR: ",
		Header: header,
		Lines:  lines,
	}.String()
}

// ErrorPrefix writes a log with a prefix to the writer provided.
func ErrorPrefix(prefix, header string, lines ...string) string {
	return cliMessage{
		Style:  DefaultStyles.Error,
		Prefix: Red(prefix) + ": ",
		Header: header,
		Lines:  lines,
	}.String()
}
