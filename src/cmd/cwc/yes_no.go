package main

import (
	"fmt"

	"lib/input"
)

func yesNoValidator(s string) error {
	switch s {
	case "", "Y", "y", "N", "n":
		return nil
	default:
		return fmt.Errorf("input must be Y or n")
	}
}

func isYes(s string) bool {
	switch s {
	case "Y", "y":
		return true
	}
	return false
}

func YesNo(prompt string, dflt bool) (bool, error) {
	d := "y"
	if dflt == false {
		d = "n"
	}
	yn, err := input.AskValidate(prompt + " [y/n]", d, yesNoValidator)
	if err != nil {
		return false, err
	}
	return isYes(yn), nil
}
