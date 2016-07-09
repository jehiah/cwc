package main

import (
	"fmt"
)

//go:generate stringer -type=Vehicle

type Vehicle int

const (
	Unknown Vehicle = 1 << iota
	Taxi
	FHV
	Other
)

func detectLicenseType(license string) Vehicle {
	if len(license) > 4 {
		return FHV
	}
	fmt.Printf("Taxi? y/n: ")
	if confirm() {
		return Taxi
	}
	return FHV
}
