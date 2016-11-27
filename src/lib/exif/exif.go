package exif

import (
	"fmt"
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type Exif struct {
	Created   time.Time
	Lat, Long float64
	// Orientation string
}

func Parse(filename string) (*Exif, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	x, err := exif.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("error parsing time from %s %s", filename, err)
	}
	e := &Exif{}
	if dt, err := x.DateTime(); err == nil {
		e.Created = dt
	}
	e.Lat, e.Long, _ = x.LatLong()
	
	// Orientation
	return e, nil

}
