// Pachage toolbox provides common utilities for use in other packages.
package toolbox

import (
	"os/user"
	"path/filepath"
	"strings"
)

// ExpandPath expands a relative or home-relative path to an absolute path.
func ExpandPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	if strings.HasPrefix(path, "~/") {
		user, err := user.Current()
		if err != nil {
			return path
		}
		return filepath.Join(user.HomeDir, path[2:])
	}
	abspath, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abspath
}
