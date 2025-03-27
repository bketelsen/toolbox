package ui

import "github.com/charmbracelet/huh"

func Prompt(title, description, placeholder string, input *string) {

	huh.NewInput().
		Title(title).
		Description(description).
		Placeholder(placeholder).
		Prompt("> ").
		Value(input).
		Run()

}
