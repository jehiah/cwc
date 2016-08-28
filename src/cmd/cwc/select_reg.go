package main

import (
	"fmt"

	"cwc/reg"

	input "github.com/tcnksm/go-input"
)

func getReg(v reg.Vehicle) (*reg.Reg, error) {
	var choices []string
	lookup := make(map[string]reg.Reg)
	for _, r := range reg.All {
		if r.Vehicle&v == 0 {
			continue
		}
		choices = append(choices, r.Description)
		lookup[r.Description] = r
	}

	selection, err := ui.Select("Violation: ", choices, &input.Options{Required: false, HideOrder: true})
	if err != nil {
		return nil, err
	}
	if selection == "" {
		return nil, nil
	}
	reg := lookup[selection]
	reg.Vehicle = v
	return &reg, nil
}

func CombineReg(a, b reg.Reg) reg.Reg {
	return reg.Reg{Type: ".", Description: fmt.Sprintf("%s and %s", a, b)}
}
