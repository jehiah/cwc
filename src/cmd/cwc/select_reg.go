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
		if r.Outdated {
			continue
		}
		if r.Vehicle&v == 0 {
			continue
		}
		var rr reg.Reg = r
		choices = append(choices, &choice{&rr})
	}

	selection, err := input.Select("Violation: ", nil, choices...)
	if err != nil {
		return nil, err
	}
	if selection == nil {
		return nil, nil
	}
	r := selection.(*choice).Reg
	r.Vehicle = v
	return r, nil
}

func CombineReg(a, b reg.Reg) reg.Reg {
	return reg.Reg{Type: ".", Description: fmt.Sprintf("%s and %s", a, b), Vehicle: a.Vehicle}
}
