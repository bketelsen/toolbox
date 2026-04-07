// copied from urfave/cli
// https://github.com/urfave/cli/blob/main/LICENSE

package toolbox

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

var isTracingOn = os.Getenv("TOOLBOX_TRACING") == "on"

// SetTracing enables or disables debug tracing output.
func SetTracing(on bool) {
	if on {
		isTracingOn = true
	} else {
		isTracingOn = false
	}
}

// Tracef outputs a formatted debug trace message to stderr with file, line, and function information.
func Tracef(format string, a ...any) {
	if !isTracingOn {
		return
	}

	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}

	pc, file, line, _ := runtime.Caller(1)
	cf := runtime.FuncForPC(pc)

	fmt.Fprintf(
		os.Stderr,
		strings.Join([]string{
			"## TOOLBOX TRACE ",
			file,
			":",
			fmt.Sprintf("%v", line),
			" ",
			fmt.Sprintf("(%s)", cf.Name()),
			" ",
			format,
		}, ""),
		a...,
	)
}
