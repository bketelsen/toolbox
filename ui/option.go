package ui

import "github.com/charmbracelet/huh"

func Option(title, description string, input *string, options []string) {

	huh.NewSelect[string]().
		Options(huh.NewOptions(options...)...).
		Title(title).
		Description(description).
		Value(input).
		Run()

}
