package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
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

func getMovieCreationTime(filePath string) (time.Time, error) {
	// Use the appropriate method to extract the date-time metadata
	// from the .mov file. This can vary depending on the operating system and available tools.
	// Here's an example command using ffprobe (FFmpeg) on Unix-like systems:
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream_tags=creation_time", "-of", "default=noprint_wrappers=1:nokey=1", filePath)
	output, err := cmd.Output()
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.Parse("2006-01-02T15:04:05.000000Z", strings.TrimSpace(string(output)))
	if err != nil {
		return t, err
	}
	nyc, _ := time.LoadLocation("America/New_York")
	return t.In(nyc), nil
}

func runNewComplaint(d db.ReadWrite) error {
	yyyymmdd, err := input.Ask("Date (YYYYMMDD) or Filename", "")
	if err != nil {
		return err
	}

	var dt time.Time
	switch {
	case strings.HasPrefix(yyyymmdd, "/"):
		ext := filepath.Ext(yyyymmdd)
		switch strings.ToLower(ext) {
		case ".jpeg", ".jpg", ".png":
			x, err := exif.ParseFile(yyyymmdd)
			if err != nil {
				return err
			}
			if x.Created.IsZero() {
				return fmt.Errorf("no timestamp found in %q", yyyymmdd)
			}
			dt = x.Created
		case ".mov":
			dt, err = getMovieCreationTime(yyyymmdd)
			if err != nil {
				return err
			}
		}
		fmt.Printf("> extracted time %s\n", dt.Format("2006/01/02 3:04pm"))
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
