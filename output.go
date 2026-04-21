package toolbox

import (
	"encoding/json"
	"fmt"
	"os"
)

// OutputJSON writes data as indented JSON to stdout if JSONOutput is true.
// Returns true if output was written, false if JSON mode is not active.
func OutputJSON(data any) bool {
	if !JSONOutput {
		return false
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		fmt.Fprintf(os.Stderr, "toolbox: JSON encode error: %v\n", err)
	}
	return true
}

// OutputJSONError writes a structured error object as JSON to stdout and
// returns a wrapped error for the caller to propagate.
func OutputJSONError(message string, err error) error {
	errOutput := map[string]any{
		"error":   true,
		"message": message,
		"details": err.Error(),
	}
	_ = OutputJSON(errOutput)
	return fmt.Errorf("%s: %w", message, err)
}
