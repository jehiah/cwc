package main

import (
	"cwc/reg"

	input "github.com/tcnksm/go-input"
)

func SelectTemplate(r reg.Reg, location string) (string, error) {
	var choices []string
	for _, s := range reg.Templates {
		if s.Code == r.Code || s.Code == "*" {
			choices = append(choices, s.Description)
		}
	}

	selection, err := ui.Select("Violation: ", choices, &input.Options{Required: true, HideOrder: true, Default: choices[0]})
	if err != nil {
		return "", err
	}

	return reg.FormatTemplate(selection, location, r.Vehicle.String(), r.String()), nil
}
