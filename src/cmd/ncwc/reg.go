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
	Short       string
	Type        string
	Vehicle     Vehicle
}

var either Vehicle = Taxi | FHV

var allReg []Reg = []Reg{
	{Code: "4-12(p)(2)", Description: "no driving in bike lane", Type: "moving", Vehicle: either},
	{Code: "4-08(e)(9)", Description: "no stopping in bike lane", Type: "parking", Vehicle: either},
	{Code: "4-11(c)(6)", Description: "no pickup or discharge of passengers in bike lane", Type: "parking", Vehicle: either, Short: "pickup/discharge in bike lane"},
	{Code: "4-08(e)(3)", Description: "no parking on sidewalks", Type: "parking", Vehicle: either},
	{Code: "4-07(b)(2)", Description: "blocking intersection and crosswalks", Type: "parking", Vehicle: either, Short: "blocking intersection/xwalk"},
	{Code: "4-05(b)(1)", Description: "no u-turns in business district", Type: "moving", Vehicle: either, Short: "no u-turns"},
	{Code: "4-12(i)", Description: "no honking in non-danger situations", Type: "parking", Vehicle: either, Short: "no honking"},
	{Code: "54-13(a)(3)(ix)", Description: "yield sign violation", Vehicle: Taxi},
	{Code: "55-13(a)(3)(ix)", Description: "yield sign violation", Vehicle: FHV},
	{Code: "54-13(a)(3)(vi)", Description: "failing to yield right of way", Vehicle: Taxi, Short: "failing to yield ROW"},
	{Code: "55-13(a)(3)(vi)", Description: "failing to yield right of way", Vehicle: FHV, Short: "failing to yield ROW"},
	{Code: "54-13(a)(3)(vii)", Description: "traffic signal violation", Vehicle: Taxi},
	{Code: "55-13(a)(3)(vii)", Description: "traffic signal violation", Vehicle: FHV},
	{Code: "54-13(a)(3)(xi)", Description: "improper passing", Vehicle: Taxi},
	{Code: "55-13(a)(3)(xi)", Description: "improper passing", Vehicle: FHV},
	{Code: "54-13(a)(3)(xii)", Description: "unsafe lane change", Vehicle: Taxi},
	{Code: "55-13(a)(3)(xii)", Description: "unsafe lane change", Vehicle: FHV},
	{Code: "NY VTL 1160(a)", Description: "no right from center lane", Type: "moving", Vehicle: either, Short: "no R from center lane"},
	{Code: "NY VTL 1160(b)", Description: "no left from center lane when both two-way streets", Type: "moving", Vehicle: either, Short: "no L from center (@ 2-way)"},
	{Code: "NY VTL 1160(c)", Description: "no left from center lane at one-way street", Type: "moving", Vehicle: either, Short: "no L from center (@ 1-way)"},
	{Code: "NY VTL 1126", Description: "no passing zone", Type: "moving", Vehicle: either},
	{Code: "NY VTL 402(b)", Description: "license plate must not be obstructed", Type: "parking", Vehicle: either, Short: "obstructed license plate"},
	{Code: "NY VTL 375(12-a)(b)(2)", Description: "no side window tint below 70%", Type: "parking", Vehicle: either, Short: "no tint below 70%"},
	{Code: "54-12(f)", Description: "threats, harassment, abuse", Vehicle: Taxi},
	{Code: "55-12(e)", Description: "threats, harassment, abuse", Vehicle: FHV},
	{Code: "54-12(g)", Description: "use or threat of physical force", Vehicle: Taxi, Short: "use/threat of physical force"},
	{Code: "55-12(f)", Description: "use or threat of physical force", Vehicle: FHV, Short: "use/threat of physical force"},
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
	{"4-12(p)(2)", "<VEHICLE> was driving in bike lane to avoid waiting in traffic in through lane, obstructing my use of bike lane <VIOLATION>. Pictures included."},
	{"4-12(p)(2)", "While near <LOCATION> I observed <VEHICLE> driving in bike lane to avoid waiting in through lane for other vehicles <VIOLATION>. Pictures included."},
	{"55-13(a)(3)(vi)", "At <LOCATION>, <VEHICLE> cut me off in the bike lane failing to yield right of way <VIOLATION>. Pictures included."},
	{"54-13(a)(3)(vii)", "At <LOCATION> I observed <VEHICLE> run red light <VIOLATION>. Pictures included. Pictures show light red and vehicle before intersection, and then vehicle proceeding through intersection on red."},
	{"55-13(a)(3)(vii)", "At <LOCATION> I observed <VEHICLE> run red light <VIOLATION>. Pictures included. Pictures show light red and vehicle before intersection, and then vehicle proceeding through intersection on red."},
	{"NY VTL 402(b)", "At <LOCATION> I observed <VEHICLE> with license plate frame obstructing view of front license plate <VIOLATION>. Pictures included show obstructed view."},
	{"NY VTL 402(b)", "At <LOCATION> I observed <VEHICLE> with license plate frame obstructing view of front license plate <VIOLATION>. NYC VTL 402(6) indicates this constitues a parking violation subject to Commission Rule 55-13(a)(1). Pictures included show obstructed view."},
	{"NY VTL 402(b)", "At <LOCATION> I observed <VEHICLE> with license plate frame obstructing view of \"T&LC\" text on rear license plate <VIOLATION>. Pictures included show obstructed view.."},
	{"NY VTL 402(b)", "At <LOCATION> I observed <VEHICLE> with license plate frame obstructing view of \"T&LC\" text on rear license plate <VIOLATION>. NYC VTL 402(6) indicates this constitues a parking violation subject to Commission Rule 55-13(a)(1). Pictures included show obstructed view."},
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
	fmt.Printf("Template: ")
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
		if r.Vehicle&v == 0 {
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
			reg.Vehicle = v
			return
		case len(selected) == 1:
			selected[0].Vehicle = v
			return selected[0]
		default:
			// join 2 together
			selected[0].Vehicle = v
			selected[1].Vehicle = v
			reg = Reg{Type: ".", Description: fmt.Sprintf("%s and %s", selected[0], selected[1])}
			selected = selected[1:]
			selected[0] = reg
		}
	}
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
		code = "Commission Rule " + code
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
	case "", "parking", "moving":
	case ".":
		return code + reg.Description
	default:
		log.Printf("unknown reg type %v", reg.Type)
		return code + reg.Description
	}

	return fmt.Sprintf("%s (%s)%s", code, reg.Description, suffix)
}
