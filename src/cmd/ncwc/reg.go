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
	Vehicle     Vehicle
}

var allReg []Reg = []Reg{
	{Code: "4-12(p)(2)", Description: "no driving in bike lane", Type: "moving", Vehicle: Taxi | FHV},
	{Code: "4-08(e)(9)", Description: "no stopping in bike lane", Type: "parking", Vehicle: Taxi | FHV},
	{Code: "4-11(c)(6)", Description: "no pickup or discharge of passengers in bike lane", Type: "parking", Vehicle: Taxi | FHV},
	{Code: "4-08(e)(3)", Description: "no parking on sidewalks", Type: "parking", Vehicle: Taxi | FHV},
	{Code: "4-07(b)(2)", Description: "blocking intersection and crosswalks", Type: "parking", Vehicle: Taxi | FHV},
	{Code: "4-05(b)(1)", Description: "no u-turns in business district", Type: "moving", Vehicle: Taxi | FHV},
	{Code: "4-12(i)", Description: "no honking in non-danger situations", Type: "parking"},
	{Code: "54-13(a)(3)(ix)", Description: "yield sign violation", Type: "-", Vehicle: Taxi},
	{Code: "55-13(a)(3)(ix)", Description: "yield sign violation", Type: "-", Vehicle: FHV},
	{Code: "54-13(a)(3)(vi)", Description: "failing to yield right of way", Type: "-", Vehicle: Taxi},
	{Code: "55-13(a)(3)(vi)", Description: "failing to yield right of way", Type: "-", Vehicle: FHV},
	{Code: "54-13(a)(3)(vii)", Description: "traffic signal violation", Type: "-", Vehicle: Taxi},
	{Code: "55-13(a)(3)(vii)", Description: "traffic signal violation", Type: "-", Vehicle: FHV},
	{Code: "54-13(a)(3)(xi)", Description: "improper passing", Type: "-", Vehicle: Taxi},
	{Code: "55-13(a)(3)(xi)", Description: "improper passing", Type: "-", Vehicle: FHV},
	{Code: "54-13(a)(3)(xii)", Description: "unsafe lane change", Type: "-", Vehicle: Taxi},
	{Code: "55-13(a)(3)(xii)", Description: "unsafe lane change", Type: "-", Vehicle: FHV},
	{Code: "NY VTL 1160(a)", Description: "no right from center lane", Type: "moving", Vehicle: Taxi | FHV},
	{Code: "NY VTL 1160(b)", Description: "no left from center lane when both two-way streets", Type: "moving", Vehicle: Taxi | FHV},
	{Code: "NY VTL 1160(c)", Description: "no left from center lane at one-way street", Type: "moving", Vehicle: Taxi | FHV},
	{Code: "NY VTL 1126", Description: "no passing zone", Type: "moving", Vehicle: Taxi | FHV},
}

type Sample struct {
	Code        string
	Description string
}

var Samples []Sample = []Sample{
	{"*", "At <LOCATION> I observed <VEHICLE> <VIOLATION>. Pictures included."},
	{"4-12(i)", "While riding bike at <LOCATION>, <VEHICLE> tried to intimidate me by honking at me <VIOLATION>. Pictures included."},
	{"4-07(b)(2)", "While biking at <LOCATION>, observed <VEHICLE> blocking crosswalk obstructing pedestrian ROW <VIOLATION>. Pictures included."},
	{"4-07(b)(2)", "While trying to cross the street at <LOCATION>, observed <VEHICLE> blocking crosswalk obstructing pedestrian ROW <VIOLATION>. Pictures included."},
	{"4-07(b)(2)", "While at <LOCATION>, observed <VEHICLE> blocking intersection and causing gridlock <VIOLATION>. Pictures included."},
	{"4-08(e)(9)", "<VEHICLE> stopped in bike lane, dangerously forcing bikers (including myself) into traffic lane <VIOLATION>. Pictures included."},
	{"4-08(e)(9)", "<VEHICLE> stopped in bike lane, obstructing my use of bike lane <VIOLATION>. Pictures included."},
	{"4-08(e)(9)", "While near <LOCATION> I observed <VEHICLE> stopped in bike lane <VIOLATION>. Pictures included."},
	{"4-12(p)(2)", "<VEHICLE> was driving in bike lane to avoid waiting in single lane obstructing my use of bike lane <VIOLATION>. Pictures included."},
	{"4-12(p)(2)", "While near <LOCATION> I observed <VEHICLE> driving in bike lane to avoid waiting in lane for other vehicles <VIOLATION>. Pictures included."},
	{"55-13(a)(3)(vi)", "At <LOCATION>, <VEHICLE> cut me off in the bike lane failing to yield right of way <VIOLATION>. Pictures included."},
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

	vehicle := fmt.Sprintf("%s Driver", reg.Vehicle)
	return o[n-1].Format(location, vehicle, reg.String())
}

func getReg(v Vehicle) (reg Reg) {
	var choices []Reg
	for _, r := range allReg {
		if r.Vehicle&v != 0 {
			continue
		}
		choices = append(choices, r)
		fmt.Printf("%d: %s\n", len(choices), r.Description)
	}

	fmt.Printf("Violation (comma separate multiple): ")
	var s string
	fmt.Scanf("%s\n", &s)
	var selected []Reg
	for _, ns := range strings.Split(s, ",") {
		n, err := strconv.Atoi(ns)
		if err != nil {
			log.Printf("%s", err)
			continue
		}
		if n < 1 || n > len(choices) {
			log.Printf("invalid option %d", n)
			continue
		}
		selected = append(selected, choices[n-1])
	}

	for {
		switch {
		case len(selected) == 0:
			break
		case len(selected) == 1:
			reg = selected[0]
			break
		default:
			// join 2 together
			selected[0].Vehicle = v
			selected[1].Vehicle = v
			reg = Reg{Type: ".", Description: fmt.Sprintf("%s and %s", selected[0], selected[1])}
			selected = selected[1:]
			selected[0] = reg
		}
	}
	reg.Vehicle = v
	return
}

func (reg Reg) String() string {
	code := reg.Code
	switch {
	case strings.HasPrefix(code, "4-"):
		code = "NYC Traffic Rule " + code
	case strings.HasPrefix(code, "54-13"):
		fallthrough
	case strings.HasPrefix(code, "55-13"):
		code = "Commision Rule " + code
	}

	var suffix string
	if !strings.Contains(code, "Commission Rule") {
		switch reg.Type {
		case "moving":
			if reg.Vehicle == FHV {
				suffix = " & Commission Rule 55-13(a)(2)"
			} else {
				suffix = " & Commission Rule 54-13(a)(2)"
			}
		case "parking":
			if reg.Vehicle == FHV {
				suffix = " & Commission Rule 55-13(a)(1)"
			} else {
				suffix = " & Commission Rule 54-13(a)(1)"
			}
		}
	}

	switch reg.Type {
	case "-", "parking", "moving":
	case ".":
		return code + reg.Description
	default:
		log.Printf("unknown reg type %v", reg.Type)
		return code + reg.Description
	}

	return fmt.Sprintf("%s (%s)%s", code, reg.Description, suffix)
}
