package main

import (
	"fmt"

	"cwc/reg"

	"lib/input"
)

type choice struct {
	*reg.Reg
}

func (c *choice) String() string {
	return c.Description
}

func getReg(v reg.Vehicle) (*reg.Reg, error) {
	var choices []interface{}
	for _, r := range reg.All {
		if r.Vehicle&v == 0 {
			continue
		}
		choices = append(choices, &choice{&r})
	}

	selection, err := input.Select("Violation", nil, choices...)
	if err != nil {
		return nil, err
	}
	if selection == nil {
		return nil, nil
	}
	r := selection.(*choice).Reg
	return r, nil
}

func CombineReg(a, b reg.Reg) reg.Reg {
	return reg.Reg{Type: ".", Description: fmt.Sprintf("%s and %s", a, b)}
}
