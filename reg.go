package main

import (
	"fmt"
	"log"
)

var parkingReg map[string]string = map[string]string{
	"NYC Traffic Rule 4-11(c)(6)": "no pickup or discharge of passengers in bike lane",
	"NYC Traffic Rule 4-08(e)(9)": "no stopping in bike lane",
	"NYC Traffic Rule 4-08(e)(3)": "no parking on sidewalks",
	"NYC Traffic Rule 4-12(i)":    "no honking in non-danger situations",
	"NYC Traffic Rule 4-07(b)(2)": "blocking intersection and crosswalks",
}

var movingReg map[string]string = map[string]string{
	"NYC Traffic Rule 4-12(p)(2)": "no driving in bike lane",
	"NYC Traffic Rule 4-05(b)(1)": "no u-turns in business district",
	"NY VTL 1160(a)":              "no right from center lane",
	"NY VTL 1160(b)":              "no left from center lane when both two-way streets",
	"NY VTL 1160(c)":              "no left from center lane at one-way street",
}

func getReg() string {
	i := 0
	var options []string
	for k, v := range parkingReg {
		fmt.Printf("%d: %s\n", i, v)
		i++
		options = append(options, k)
	}
	for k, v := range movingReg {
		fmt.Printf("%d: %s\n", i, v)
		i++
		options = append(options, k)
	}

	fmt.Printf("Violation: ")
	var n int
	fmt.Scanf("%d\n", &n)
	if n > len(options) {
		log.Printf("invalid option %d", n)
		return ""
	}
	return options[n]
}

func formatReg(reg string, fhv bool) string {
	var tlcReg string
	if txt, ok := movingReg[reg]; ok {
		if fhv {
			tlcReg = "Commission Rule 55-13(a)(2)"
		} else {
			tlcReg = "Commission Rule 54-13(a)(2)"
		}
		return fmt.Sprintf("in violation of %s (%s) & %s", reg, txt, tlcReg)
	}
	if txt, ok := parkingReg[reg]; ok {
		if fhv {
			tlcReg = "Commission Rule 55-13(a)(1)"
		} else {
			tlcReg = "Commission Rule 54-13(a)(1)"
		}
		return fmt.Sprintf("in violation of %s (%s) & %s", reg, txt, tlcReg)
	}
	return reg
}
