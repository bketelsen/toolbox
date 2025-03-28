package ui

import (
	"io"
	"os"
)

var defaultLinePrefix = "*"

type ConsolePrinter struct {
	stdout     io.Writer
	stderr     io.Writer
	linePrefix string
}

var Console *ConsolePrinter = New(os.Stdout, os.Stderr)

func New(out io.Writer, err io.Writer) *ConsolePrinter {
	return &ConsolePrinter{
		stdout:     out,
		stderr:     err,
		linePrefix: defaultLinePrefix,
	}
}

func (c *ConsolePrinter) SetStderr(err io.Writer) {
	c.stderr = err
}

func (c *ConsolePrinter) SetStdout(out io.Writer) {
	c.stdout = out
}

func (c *ConsolePrinter) SetLinePrefix(prefix string) {
	c.linePrefix = prefix
}

func SetStdout(out io.Writer) {
	Console.SetStdout(out)
}

func SetStderr(err io.Writer) {
	Console.SetStderr(err)
}

func SetLinePrefix(prefix string) {
	Console.SetLinePrefix(prefix)
}

func (c ConsolePrinter) Warn(header string, lines ...string) {
	_, _ = c.stderr.Write([]byte(warn(c.linePrefix, header, lines...)))
}

func (c ConsolePrinter) WarnPrefix(prefix, header string, lines ...string) {
	_, _ = c.stderr.Write([]byte(warnPrefix(c.linePrefix, prefix, header, lines...)))
}

func (c ConsolePrinter) Info(header string, lines ...string) {
	_, _ = c.stderr.Write([]byte(info(c.linePrefix, header, lines...)))
}

func (c ConsolePrinter) InfoPrefix(prefix, header string, lines ...string) {
	_, _ = c.stderr.Write([]byte(infoPrefix(c.linePrefix, prefix, header, lines...)))
}

func (c ConsolePrinter) Success(header string, lines ...string) {
	_, _ = c.stderr.Write([]byte(success(c.linePrefix, header, lines...)))
}

func (c ConsolePrinter) SuccessPrefix(prefix, header string, lines ...string) {
	_, _ = c.stderr.Write([]byte(successPrefix(c.linePrefix, prefix, header, lines...)))
}

func (c ConsolePrinter) Error(header string, lines ...string) {
	_, _ = c.stderr.Write([]byte(printerror(c.linePrefix, header, lines...)))
}

func (c ConsolePrinter) ErrorPrefix(prefix, header string, lines ...string) {
	_, _ = c.stderr.Write([]byte(errorPrefix(c.linePrefix, prefix, header, lines...)))
}

func Warn(header string, lines ...string) {
	Console.Warn(header, lines...)
}

func WarnPrefix(prefix, header string, lines ...string) {
	Console.WarnPrefix(prefix, header, lines...)
}

func Info(header string, lines ...string) {
	Console.Info(header, lines...)
}

func InfoPrefix(prefix, header string, lines ...string) {
	Console.InfoPrefix(prefix, header, lines...)
}

func Success(header string, lines ...string) {
	Console.Success(header, lines...)
}

func SuccessPrefix(prefix, header string, lines ...string) {
	Console.SuccessPrefix(prefix, header, lines...)
}

func Error(header string, lines ...string) {
	Console.Error(header, lines...)
}

func ErrorPrefix(prefix, header string, lines ...string) {
	Console.ErrorPrefix(prefix, header, lines...)
}
