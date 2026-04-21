package ui

import "charm.land/huh/v2"

// Confirm displays an interactive yes/no confirmation and writes the result to confirm.
// Returns an error if the prompt is cancelled or the terminal is unavailable.
func Confirm(title, description string, confirm *bool) error {

	return huh.NewConfirm().
		Title(title).
		Description(description).
		Value(confirm).
		Run()

}
