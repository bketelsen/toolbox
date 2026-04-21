package toolbox

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestRegisterFlags(t *testing.T) {
	app := &App{}
	cmd := &cobra.Command{Use: "test"}

	app.registerFlags(cmd)

	flags := []struct {
		name      string
		shorthand string
	}{
		{"json", ""},
		{"verbose", "v"},
		{"dry-run", "n"},
		{"silent", "s"},
		{"log-file", ""},
	}

	for _, f := range flags {
		pf := cmd.PersistentFlags().Lookup(f.name)
		if pf == nil {
			t.Errorf("flag --%s not registered", f.name)
			continue
		}
		if f.shorthand != "" && pf.Shorthand != f.shorthand {
			t.Errorf("flag --%s shorthand = %q, want %q", f.name, pf.Shorthand, f.shorthand)
		}
	}
}

func TestBindViper(t *testing.T) {
	app := &App{}
	cmd := &cobra.Command{Use: "test"}
	app.registerFlags(cmd)

	err := BindViper(cmd)
	if err != nil {
		t.Fatalf("BindViper() error = %v", err)
	}
}

