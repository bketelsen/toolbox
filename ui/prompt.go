package ui

import "charm.land/huh/v2"

func Prompt(title, description, placeholder string, input *string) {

	huh.NewInput().
		Title(title).
		Description(description).
		Placeholder(placeholder).
		Prompt("> ").
		Value(input).
		Run()

}
