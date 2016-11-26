package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"cwc/db"
	"cwc/reg"

	"lib/exif"
	"lib/input"
)

func newComplaint() error {

	yyyymmdd, err := input.Ask("Date (YYYYMMDD) or Filename", "")
	if err != nil {
		return err
	}

	var dt time.Time
	switch {
	case strings.HasPrefix(yyyymmdd, "/"):
		x, err := exif.Parse(yyyymmdd)
		if err != nil {
			return err
		}
		if x.Created.IsZero() {
			return fmt.Errorf("no timestamp found in %q", yyyymmdd)
		}
		dt = x.Created
		fmt.Printf("> using EXIF time %s\n", dt.Format("2006/01/02 3:04pm"))
	case yyyymmdd == "":
		yyyymmdd = time.Now().Format("20060102")
		fmt.Printf(" > using %s\n", yyyymmdd)
		fallthrough
	default:
		hhmm, err := input.Ask("Time (HHMM)", "")
		if err != nil {
			return err
		}
		dt, err = time.Parse("20060102 1504", fmt.Sprintf("%s %s", yyyymmdd, hhmm))
		if err != nil {
			return err
		}
	}

	license, err := input.Ask("License Plate", "")
	if err != nil {
		return err
	}

	vehicle := reg.FHV
	if reg.PossibleTaxi(license) {
		if err != nil {
			return err
		}
		yn, err := YesNo("Taxi", true)
		if err != nil {
			return err
		}
		if yn {
			vehicle = reg.Taxi
		}
	}

	where, err := input.Ask("Where", "")
	if err != nil {
		return err
	}

	complaint, err := db.Default.New(dt, license)
	if err != nil {
		return err
	}
	fmt.Printf("> creating %s\n", db.Default.FullPath(complaint))
	f, err := db.Default.Create(complaint)
	if err != nil {
		return err
	}

	fmt.Fprintf(f, "%s %s %s %s\n", dt.Format("2006/01/02 3:04pm"), vehicle, license, where)

	for {
		r, err := getReg(vehicle)
		if err != nil {
			return err
		}
		if r == nil {
			return errors.New("no regulation selected")
		}
		body, err := SelectTemplate(*r, where)
		if err != nil {
			return err
		}
		fmt.Fprintf(f, "\n%s\n", body)
		yn, err := YesNo("Another Violation", false)
		if err != nil {
			return err
		}
		if !yn {
			break
		}
	}

	f.Close()

	var url string
	if vehicle == reg.FHV {
		url = "https://www1.nyc.gov/apps/311universalintake/form.htm?serviceName=TLC+FHV+Driver+Unsafe+Driving"
	} else {
		url = "https://www1.nyc.gov/apps/311universalintake/form.htm?serviceName=TLC+Taxi+Driver+Unsafe+Driving+Non-Passenger"
	}
	err = exec.Command("/usr/bin/open", "-a", "/Applications/Google Chrome.app/", url).Run()
	if err != nil {
		return err
	}
	fmt.Printf("> opening %s\n", url)
	db.Default.Edit(complaint)
	db.Default.ShowInFinder(complaint)
	fmt.Printf("> done\n")
	return nil
}
