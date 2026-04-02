package ui

import "charm.land/huh/v2"

func Confirm(title, description string, confirm *bool) {

	huh.NewConfirm().
		Title(title).
		Description(description).
		Value(confirm).
		Run()

}
