package main

import (
	"fmt"
	"log"
	"strings"
)

type Reg struct {
	Code        string
	Description string
	Type        string
	FHV         bool
}

var allReg []Reg = []Reg{
	{Code: "4-11(c)(6)", Description: "no pickup or discharge of passengers in bike lane", Type: "parking"},
	{Code: "4-08(e)(9)", Description: "no stopping in bike lane", Type: "parking"},
	{Code: "4-08(e)(3)", Description: "no parking on sidewalks", Type: "parking"},
	{Code: "4-12(i)", Description: "no honking in non-danger situations", Type: "parking"},
	{Code: "4-07(b)(2)", Description: "blocking intersection and crosswalks", Type: "parking"},
	{Code: "4-12(p)(2)", Description: "no driving in bike lane", Type: "moving"},
	{Code: "4-05(b)(1)", Description: "no u-turns in business district", Type: "moving"},
	{Code: "NY VTL 1160(a)", Description: "no right from center lane", Type: "moving"},
	{Code: "NY VTL 1160(b)", Description: "no left from center lane when both two-way streets", Type: "moving"},
	{Code: "NY VTL 1160(c)", Description: "no left from center lane at one-way street", Type: "moving"},
}

type Sample struct {
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

func (reg Reg) String() string {
	var tlcReg string
	switch reg.Type {
	case "moving":
		if reg.FHV {
			tlcReg = "Commission Rule 55-13(a)(2)"
		} else {
			tlcReg = "Commission Rule 54-13(a)(2)"
		}
	case "parking":
		if reg.FHV {
			tlcReg = "Commission Rule 55-13(a)(1)"
		} else {
			tlcReg = "Commission Rule 54-13(a)(1)"
		}
	default:
		log.Printf("unknown reg type %v", reg.Type)
		return ""
	}

	code := reg.Code
	if !strings.HasPrefix(code, "NY VTL ") {
		code = fmt.Sprintf("NYC Traffic Rule %s", code)
	}

	return fmt.Sprintf("in violation of %s (%s) & %s", code, reg.Description, tlcReg)
}

func expandReg(r string) string {
	switch {
	case strings.Contains(r, "driving in bike lane"):
		return fmt.Sprintf("driving in bike lane, dangerously forcing bikers (including myself) into traffic lane, %s", r)
	case strings.Contains(r, "bike lane"):
		return fmt.Sprintf("stopped in bike lane, dangerously forcing bikers (including myself) into traffic lane, %s", r)
	}
	return r
}
