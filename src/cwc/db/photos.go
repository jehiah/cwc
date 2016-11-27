package db

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"lib/exif"
)

type Photo struct {
	Name      string
	Filename  string
	Submitted bool
	Created   time.Time
	Lat, Long float64
	// Orientation
}

func (f *FullComplaint) ParsePhotos() {
	if len(f.PhotoDetails) > 0 && len(f.Photos) > 0 {
		// idempotent
		return
	}

	// build up index of photos that were noted as submitted
	submitted := make(map[string]bool)
	skip := make(map[string]bool)
	if len(f.Lines) >= 3 && strings.HasPrefix(f.Lines[2], "photos") {
		for _, s := range strings.Split(strings.Replace(f.Lines[2], " ", ",", -1), ",") {
			if s == "photos" {
				continue
			}
			for _, file := range f.Photos {
				prefix := file[:len(file)-len(filepath.Ext(file))]
				if strings.HasSuffix(prefix, s) {
					submitted[file] = true
					break
				}
			}
		}
	}

	// catch 'copy' records
	// side effect of workflow for resized entries
	for _, file := range f.Photos {
		if strings.Contains(file, " copy.") {
			real := strings.Replace(file, " copy.", ".", 1)
			submitted[real] = true
			skip[file] = true
		}
	}

	for _, file := range f.Photos {
		if skip[file] {
			continue
		}
		p := Photo{
			Name:      file,
			Filename:  filepath.Join(f.BasePath, file),
			Submitted: submitted[file],
		}
		x, err := exif.Parse(p.Filename)
		if err != nil {
			// TODO: use file created timestamp
			log.Printf("%s", err)
		} else {
			p.Created = x.Created
			p.Lat, p.Long = x.Lat, x.Long
		}
		f.PhotoDetails = append(f.PhotoDetails, &p)
	}
	
	// TODO: sort by timestamps
	return
}

func (f *FullComplaint) HasGPSInfo() bool {
	f.ParsePhotos()
	for _, p := range f.PhotoDetails {
		if p.Lat != 0 && p.Long != 0 {
			return true
		}
	}
	return false
}
