package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"cwc/db"
	"cwc/reg"

	input "github.com/tcnksm/go-input"
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

func newComplaint() error {

	yyyymmdd, err := ui.Ask("Date (YYYYMMDD) or Filename:", &input.Options{Required: true, Loop: true, HideOrder: true})
	if err != nil {
		return err
	}

	var dt time.Time
	switch {
	case strings.HasPrefix(yyyymmdd, "/"):
		dt, err = getExifDateTime(yyyymmdd)
		if err != nil {
			return err
		}
		fmt.Printf("> using EXIF time %s\n", dt.Format("2006/01/02 3:04pm"))
	case yyyymmdd == "":
		yyyymmdd = time.Now().Format("20060102")
		fmt.Printf(" > using %s\n", yyyymmdd)
		fallthrough
	default:
		hhmm, err := ui.Ask("Time (HHMM): ", &input.Options{Required: true, Loop: true, HideOrder: true})
		if err != nil {
			return err
		}
		dt, err = time.Parse("20060102 1504", fmt.Sprintf("%s %s", yyyymmdd, hhmm))
		if err != nil {
			return err
		}
	}

	license, err := ui.Ask("License Plate: ", &input.Options{Required: true, Loop: true, HideOrder: true})
	if err != nil {
		return err
	}

	vehicle := reg.FHV
	if reg.PossibleTaxi(license) {
		if err != nil {
			return err
		}
		yn, err := ui.Ask("Taxi: ", &input.Options{Required: true, Loop: true, HideOrder: true, ValidateFunc: yesNoValidator})
		if err != nil {
			return err
		}
		if isYes(yn) {
			vehicle = reg.Taxi
		}
	}

	where, err := ui.Ask("Where? ", &input.Options{Required: true, Loop: true, HideOrder: true})
	if err != nil {
		return err
	}
	where = strings.TrimSpace(where)

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

	var r *reg.Reg
	for {
		n, err := getReg(vehicle)
		if err != nil {
			return err
		}
		if r == nil && n == nil {
			return errors.New("no regulation selected")
		}
		if r == nil {
			break
		}
		if n == nil {
			r = n
		} else {
			rr := CombineReg(*r, *n)
			r = &rr
		}
	}

	body, err := SelectTemplate(*r, where)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, "\n%s\n", body)
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
