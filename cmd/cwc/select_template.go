package main

import (
	"github.com/jehiah/cwc/input"
	"github.com/jehiah/cwc/internal/reg"
)

func SelectTemplate(r reg.Reg, location, license string) (string, error) {
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

	return reg.FormatTemplate(selection, location, r.Vehicle.String(), license, r.String()), nil
}
