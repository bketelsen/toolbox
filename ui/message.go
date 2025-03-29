// https://github.com/coder/coder/blob/main/LICENSE
// Extracted and modified from github.com/coder/coder
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	LinePrefixDefault     = "|"
	LinePrefixBullet      = "•"
	LinePrefixCheck       = "✓"
	LinePrefixCross       = "✗"
	LinePrefixWarning     = "⚠"
	LinePrefixError       = "✖"
	LinePrefixSuccess     = "✔"
	LinePrefixQuestion    = "?"
	LinePrefixDoubleAngle = "»"
)

// cliMessage provides a human-readable message for CLI errors and messages.
type cliMessage struct {
	Style      lipgloss.Style
	Header     string
	Prefix     string
	Lines      []string
	LinePrefix string
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
		_, _ = fmt.Fprintf(&str, "  %s %s\r\n", m.Style.Render(m.LinePrefix), line)
	}
	return str.String()
}

// Warn writes a log to the writer provided.
func warn(linePrefix, header string, lines ...string) string {
	return cliMessage{
		LinePrefix: linePrefix,
		Style:      DefaultStyles.Warn,
		Prefix:     "WARNING: ",
		Header:     header,
		Lines:      lines,
	}.String()
}

func warnPrefix(linePrefix, prefix, header string, lines ...string) string {
	return cliMessage{
		LinePrefix: linePrefix,
		Style:      DefaultStyles.Warn,
		Prefix:     Yellow(prefix) + ": ",
		Header:     header,
		Lines:      lines,
	}.String()
}

func info(linePrefix, header string, lines ...string) string {
	return cliMessage{
		LinePrefix: linePrefix,
		Header:     header,
		Lines:      lines,
	}.String()
}

func success(linePrefix, header string, lines ...string) string {
	return cliMessage{
		LinePrefix: linePrefix,
		Header:     Green(header),
		Lines:      lines,
	}.String()
}

func infoPrefix(linePrefix, prefix, header string, lines ...string) string {
	return cliMessage{
		LinePrefix: linePrefix,
		Header:     header,
		Prefix:     prefix + ": ",
		Lines:      lines,
	}.String()
}

func successPrefix(linePrefix, prefix, header string, lines ...string) string {
	return cliMessage{
		Header:     header,
		LinePrefix: linePrefix,

		Prefix: Green(prefix) + ": ",
		Lines:  lines,
	}.String()
}

func printerror(linePrefix, header string, lines ...string) string {
	return cliMessage{
		LinePrefix: linePrefix,

		Style:  DefaultStyles.Error,
		Prefix: "ERROR: ",
		Header: header,
		Lines:  lines,
	}.String()
}

func errorPrefix(linePrefix, prefix, header string, lines ...string) string {
	return cliMessage{
		LinePrefix: linePrefix,
		Style:      DefaultStyles.Error,
		Prefix:     Red(prefix) + ": ",
		Header:     header,
		Lines:      lines,
	}.String()
}
