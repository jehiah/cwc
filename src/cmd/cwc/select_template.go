package main

import (
	"cwc/reg"

	"lib/input"
)

func SelectTemplate(r reg.Reg, location string) (string, error) {
	var choices []string
	for _, s := range reg.Templates {
		if s.Code == r.Code || s.Code == "*" {
			choices = append(choices, s.Description)
		}
	}

	selection, err := input.SelectString("Violation", choices[0], choices...)
	if err != nil {
		return "", err
	}

	return reg.FormatTemplate(selection, location, r.Vehicle.String(), r.String()), nil
}
