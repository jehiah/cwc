package main

import (
	"fmt"

	"cwc/reg"
)

func listRegulations(short bool) {
	for _, r := range reg.All {
		desc := r.Description
		if short && r.Short != "" {
			desc = r.Short
		}
		fmt.Printf("%s,%s\n", r.Code, desc)
	}
}
