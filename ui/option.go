package ui

import "charm.land/huh/v2"

// Option displays an interactive selection menu and writes the chosen value to input.
// Returns an error if the prompt is cancelled or the terminal is unavailable.
func Option(title, description string, input *string, options []string) error {

	return huh.NewSelect[string]().
		Options(huh.NewOptions(options...)...).
		Title(title).
		Description(description).
		Value(input).
		Run()

}
