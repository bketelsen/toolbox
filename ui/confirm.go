package ui

import "github.com/charmbracelet/huh"

func Confirm(title, description string, confirm *bool) {

	huh.NewConfirm().
		Title(title).
		Description(description).
		Value(confirm).
		Run()

}
