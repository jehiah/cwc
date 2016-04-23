package main

import (
	"fmt"
	"log"
)

type Reg struct {
	Code        string
	Description string
	Type        string
}

var allReg []Reg = []Reg{
	{"NYC Traffic Rule 4-11(c)(6)", "no pickup or discharge of passengers in bike lane", "parking"},
	{"NYC Traffic Rule 4-08(e)(9)", "no stopping in bike lane", "parking"},
	{"NYC Traffic Rule 4-08(e)(3)", "no parking on sidewalks", "parking"},
	{"NYC Traffic Rule 4-12(i)", "no honking in non-danger situations", "parking"},
	{"NYC Traffic Rule 4-07(b)(2)", "blocking intersection and crosswalks", "parking"},
	{"NYC Traffic Rule 4-12(p)(2)", "no driving in bike lane", "moving"},
	{"NYC Traffic Rule 4-05(b)(1)", "no u-turns in business district", "moving"},
	{"NY VTL 1160(a)", "no right from center lane", "moving"},
	{"NY VTL 1160(b)", "no left from center lane when both two-way streets", "moving"},
	{"NY VTL 1160(c)", "no left from center lane at one-way street", "moving"},
}

func getReg() Reg {
	for i, r := range allReg {
		fmt.Printf("%d: %s\n", i+1, r.Description)
	}

	fmt.Printf("Violation: ")
	var n int
	fmt.Scanf("%d\n", &n)
	if n < 1 || n > len(allReg) {
		log.Printf("invalid option %d", n)
		return Reg{}
	}
	return allReg[n-1]
}

func formatReg(reg Reg, fhv bool) string {
	var tlcReg string
	switch reg.Type {
	case "moving":
		if fhv {
			tlcReg = "Commission Rule 55-13(a)(2)"
		} else {
			tlcReg = "Commission Rule 54-13(a)(2)"
		}
	case "parking":
		if fhv {
			tlcReg = "Commission Rule 55-13(a)(1)"
		} else {
			tlcReg = "Commission Rule 54-13(a)(1)"
		}
	default:
		log.Printf("unknown reg type %v", reg.Type)
		return ""
	}
	return fmt.Sprintf("in violation of %s (%s) & %s", reg.Code, reg.Description, tlcReg)
}
