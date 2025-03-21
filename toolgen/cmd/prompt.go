package cmd

import "github.com/charmbracelet/huh"

func prompt(title, description, placeholder string, input *string) {

	huh.NewInput().
		Title(title).
		Description(description).
		Placeholder(placeholder).
		Prompt("> ").
		Value(input).
		Run()

}
