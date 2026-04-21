package ui

import "charm.land/huh/v2"

// Prompt displays an interactive text input and writes the result to input.
// Returns an error if the prompt is cancelled or the terminal is unavailable.
func Prompt(title, description, placeholder string, input *string) error {

	return huh.NewInput().
		Title(title).
		Description(description).
		Placeholder(placeholder).
		Prompt("> ").
		Value(input).
		Run()

}
