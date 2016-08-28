package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

func getExifDateTime(filename string) (time.Time, error) {
	f, err := os.Open(filename)
	if err != nil {
		return time.Time{}, err
	}
	x, err := exif.Decode(f)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing time from %s %s", filename, err)
	}
	dt, err := x.DateTime()
	if err != nil {
		return time.Time{}, fmt.Errorf("no EXIF date time in %s %s", filename, err)
	}
	return dt, nil

}
