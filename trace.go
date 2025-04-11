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

func SetTracing(on bool) {
	if on {
		isTracingOn = true
	} else {
		isTracingOn = false
	}
}
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
