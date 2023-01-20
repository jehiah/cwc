package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/jehiah/cwc/db"
	"github.com/jehiah/cwc/exif"
	"github.com/jehiah/cwc/input"
	"github.com/jehiah/cwc/internal/complaint"
	"github.com/jehiah/cwc/internal/reg"
	"github.com/spf13/cobra"
)

func newComplaint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "New Complaint",
		Run: func(cmd *cobra.Command, args []string) {
			err := runNewComplaint(loadDB(cmd.Flags().GetString("db")))
			if err != nil {
				log.Fatal(err.Error())
			}
		},
	}
	cmd.Flags().String("db", string(db.Default), "DB path")
	return cmd
}

func runNewComplaint(d db.ReadWrite) error {
	yyyymmdd, err := input.Ask("Date (YYYYMMDD) or Filename", "")
	if err != nil {
		return err
	}

	var dt time.Time
	switch {
	case strings.HasPrefix(yyyymmdd, "/"):
		x, err := exif.ParseFile(yyyymmdd)
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

	latest := &complaint.FullComplaint{}
	if c, err := d.Latest(); err == nil {
		if delta := c.Time().Sub(dt); delta > (-1*time.Hour) && delta < time.Hour {
			latest, _ = d.FullComplaint(c)
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

	where, err := input.Ask("Where", latest.Location)
	if err != nil {
		return err
	}

	c := complaint.New(dt, license)

	fmt.Printf("> creating %s\n", d.FullPath(c))
	f, err := d.Create(c)
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
		body, err := SelectTemplate(*r, where, license)
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
		url = "https://portal.311.nyc.gov/article/?kanumber=KA-01244"
		// url = "https://portal.311.nyc.gov/sr-step/?id=60828b57-284c-ed11-97b2-2818785c4829&stepid=3f227a9c-fb0d-e811-8127-1458d04d2538"
	} else {
		url = "https://portal.311.nyc.gov/article/?kanumber=KA-01241"
		// url = "https://www1.nyc.gov/apps/311universalintake/form.htm?serviceName=TLC+Taxi+Driver+Unsafe+Driving+Non-Passenger"
	}
	err = exec.Command("/usr/bin/open", "-a", "/Applications/Google Chrome.app/", url).Run()
	if err != nil {
		return err
	}
	if id, ok := d.(db.Interactive); ok {
		fmt.Printf("> opening %s\n", url)
		id.ShowInEditor(c)
		id.ShowInFinder(c)
	}
	fmt.Printf("> done\n")
	return nil
}
