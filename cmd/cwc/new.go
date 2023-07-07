package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jehiah/cwc/db"
	"github.com/jehiah/cwc/exif"
	"github.com/jehiah/cwc/input"
	"github.com/jehiah/cwc/internal/complaint"
	"github.com/jehiah/cwc/internal/reg"
	"github.com/jehiah/nycgeosearch"
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
	ctx := context.Background()
	yyyymmdd, err := input.Ask("Date (YYYYMMDD) or Filename", "")
	if err != nil {
		return err
	}

	var dt time.Time
	var x exif.Exif
	switch {
	case strings.HasPrefix(yyyymmdd, "/"):
		x, err = exif.ParseImageOrVideo(yyyymmdd)
		if err != nil {
			return err
		}
		if x.Created.IsZero() {
			return fmt.Errorf("no timestamp found in %q", yyyymmdd)
		}
		dt = x.Created
		fmt.Printf("> extracted time %s\n", dt.Format("2006/01/02 3:04pm"))
	case yyyymmdd == "":
		yyyymmdd = time.Now().Format("20060102")
		fmt.Printf(" > using %s\n", yyyymmdd)
		fallthrough
	default:
		hhmm, err := input.Ask("Time (HHMM lcl)", "")
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

	var nearestAddress string
	if x.Lat != 0 {
		geo, err := nycgeosearch.PlanningLabs.ReverseGeocode(ctx, nycgeosearch.Location{
			Lat: x.Lat,
			Lng: x.Long,
		}, nycgeosearch.Options{Size: 1})
		if err != nil {
			log.Printf("error doing reverse geo lookup %s", err)
		} else {
			if len(geo.Features) >= 1 {
				nearestAddress = geo.Features[0].PropertyMustString("label")
				fmt.Printf("> Nearest Address: %s\n", nearestAddress)
			}
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

	// TODO: convert HEIC to jpeg

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

	if nearestAddress != "" {
		fmt.Fprintf(f, "Address: %s\n", nearestAddress)
	}
	if x.Lat != 0 {
		fmt.Fprintf(f, "[ll:%f,%f]\n", x.Lat, x.Long)
	}

	f.Close()

	if id, ok := d.(db.Interactive); ok {
		id.ShowInEditor(c)
		id.ShowInFinder(c)
	}

	yn, err := YesNo("Submit", true)
	if yn {
		fc, err := d.FullComplaint(c)
		if err != nil {
			return err
		}
		err = Submit(fc)
		if err != nil {
			return err
		}
	}

	fmt.Printf("> done\n")
	return nil
}
