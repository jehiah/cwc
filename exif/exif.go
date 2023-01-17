package exif

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type Exif struct {
	Created   time.Time
	Lat, Long float64

	// Orientation string
	ExifRotation float64
	ExifFlip     bool
}

func ParseFile(filename string) (*Exif, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	e, err := Parse(f)
	if err != nil {
		return nil, fmt.Errorf("error parsing time from %s %w", filename, err)
	}
	return e, nil
}

func Parse(f io.Reader) (*Exif, error) {
	x, err := exif.Decode(f)
	if err != nil {
		return nil, err
	}
	e := &Exif{}
	if dt, err := x.DateTime(); err == nil {
		e.Created = dt
	}
	e.Lat, e.Long, _ = x.LatLong()

	// Orientation
	o, err := x.Get(exif.Orientation)
	if err == nil {
		n, err := o.Int(0)
		if err == nil {
			e.ExifRotation, e.ExifFlip = calculateRotationFilp(n)
		}
	}
	return e, nil

}

func calculateRotationFilp(orientation int) (rotate float64, flip bool) {
	// from https://github.com/h2non/bimg/blob/master/resize.go#L457
	switch orientation {
	case 6:
		rotate = 90
	case 3:
		rotate = 180
	case 8:
		rotate = 270
	case 2:
		flip = true
		// flip 1
	case 7:
		flip = true
		rotate = 90
		// flip 6
	case 4:
		flip = true
		rotate = 180
		// flip 3
	case 5:
		flip = true
		rotate = 270
		// flip 8
	case 1:
	default:
		log.Printf("orientation %v", orientation)
	}
	return
}
