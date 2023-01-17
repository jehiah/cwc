package complaint

import (
	"time"
)

type Photo struct {
	Name      string
	Filename  string
	Submitted bool
	Created   time.Time
	Lat, Long float64
	Size      int64
	// Orientation
}
