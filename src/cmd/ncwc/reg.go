package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Reg struct {
	Code        string
	Description string
	Type        string
	FHV         bool
}

var allReg []Reg = []Reg{
	{Code: "4-12(p)(2)", Description: "no driving in bike lane", Type: "moving"},
	{Code: "4-08(e)(9)", Description: "no stopping in bike lane", Type: "parking"},
	{Code: "4-11(c)(6)", Description: "no pickup or discharge of passengers in bike lane", Type: "parking"},
	{Code: "4-08(e)(3)", Description: "no parking on sidewalks", Type: "parking"},
	{Code: "4-07(b)(2)", Description: "blocking intersection and crosswalks", Type: "parking"},
	{Code: "4-05(b)(1)", Description: "no u-turns in business district", Type: "moving"},
	{Code: "54-13(a)(3)(xi)", Description: "improper passing", Type: "-"},
	{Code: "4-12(i)", Description: "no honking in non-danger situations", Type: "parking"},
	{Code: "NY VTL 1160(a)", Description: "no right from center lane", Type: "moving"},
	{Code: "NY VTL 1160(b)", Description: "no left from center lane when both two-way streets", Type: "moving"},
	{Code: "NY VTL 1160(c)", Description: "no left from center lane at one-way street", Type: "moving"},
	{Code: "NY VTL 1126", Description: "no passing zone", Type: "moving"},
}

type Sample struct {
	Code        string
	Description string
}

var Samples []Sample = []Sample{
	{"*", "At <LOCATION> I observed <VEHICLE> <VIOLATION>. Pictures included."},
	{"4-07(b)(2)", "While biking at <LOCATION>, observed <VEHICLE> blocking crosswalk obstructing pedestrian ROW <VIOLATION>. Pictures included."},
	{"4-07(b)(2)", "While trying to cross the street at <LOCATION>, observed <VEHICLE> blocking crosswalk obstructing pedestrian ROW <VIOLATION>. Pictures included."},
	{"4-07(b)(2)", "While at <LOCATION>, observed <VEHICLE> blocking intersection and causing gridlock <VIOLATION>. Pictures included."},
	{"55-13(a)(3)(vi)", "At <LOCATION>, <VEHICLE> cut me off in the bike lane failing to yield right of way <VIOLATION>. Pictures included."},
	{"4-12(i)", "While riding bike at <LOCATION>, <VEHICLE> tried to intimidate me by honking at me <VIOLATION>. Pictures included."},
	{"4-08(e)(9)", "<VEHICLE> stopped in bike lane, dangerously forcing bikers (including myself) into traffic lane <VIOLATION>. Pictures included."},
	{"4-08(e)(9)", "<VEHICLE> stopped in bike lane, obstructing my use of bike lane <VIOLATION>. Pictures included."},
	{"4-08(e)(9)", "While near <LOCATION> I observed <VEHICLE> stopped in bike lane <VIOLATION>. Pictures included."},
	{"4-12(p)(2)", "<VEHICLE> was driving in bike lane to avoid waiting in single lane obstructing my use of bike lane <VIOLATION>. Pictures included."},
	{"4-12(p)(2)", "While near <LOCATION> I observed <VEHICLE> driving in bike lane to avoid waiting in lane for other vehicles <VIOLATION>. Pictures included."},
}

func (s Sample) Format(location, vehicle, violation string) string {
	t := s.Description
	t = strings.Replace(t, "<LOCATION>", location, -1)
	t = strings.Replace(t, "<VEHICLE>", vehicle, -1)
	return strings.Replace(t, "<VIOLATION>", fmt.Sprintf("in violation of %s", violation), -1)
}

func SelectSample(reg Reg, location string) string {
	var o []Sample
	for _, s := range Samples {
		if s.Code == reg.Code || s.Code == "*" {
			o = append(o, s)
			fmt.Printf("%d: %s\n", len(o), s.Description)
		}
	}

	var n int
	fmt.Scanf("%d\n", &n)
	if n < 1 || n > len(o) {
		log.Printf("invalid option %d", n)
		return ""
	}

	vehicle := fmt.Sprintf("%s Driver", fhvStr(reg.FHV))
	return o[n-1].Format(location, vehicle, reg.String())
}

func getReg(fhv bool) (reg Reg) {
	for i, r := range allReg {
		fmt.Printf("%d: %s\n", i+1, r.Description)
	}

	fmt.Printf("Violation: ")
	var s string
	fmt.Scanf("%s\n", &s)
	var regs []Reg
	for _, ns := range strings.Split(s, ",") {
		n, err := strconv.Atoi(ns)
		if err != nil {
			log.Printf("%s", err)
			continue
		}
		if n < 1 || n > len(allReg) {
			log.Printf("invalid option %d", n)
			continue
		}
		regs = append(regs, allReg[n-1])
	}

	switch {
	case len(regs) == 0:
	case len(regs) == 1:
		reg = regs[0]
	default:
		// join 2 together
		regs[0].FHV = fhv
		regs[1].FHV = fhv
		reg = Reg{Type: ".", Description: fmt.Sprintf("%s and %s", regs[0], regs[1])}
	}
	reg.FHV = fhv
	return
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
	case "-":
	case ".":
		return reg.Code + reg.Description
	default:
		log.Printf("unknown reg type %v", reg.Type)
		return reg.Code + reg.Description
	}

	code := reg.Code
	if !strings.HasPrefix(code, "NY VTL ") {
		code = fmt.Sprintf("NYC Traffic Rule %s", code)
	}

	if tlcReg != "" {
		tlcReg = fmt.Sprintf(" & %s", tlcReg)
	}

	return fmt.Sprintf("%s (%s)%s", code, reg.Description, tlcReg)
}
