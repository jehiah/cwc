package db

import (
	"log"
	"os"
	"path/filepath"
	"sort"
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
	Size      int64
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
			if s == "photos" || s == "" {
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
	// side effect of workflow for resized/annotated entries
	for _, file := range f.Photos {
		if strings.Contains(file, " copy.") {
			real := strings.Replace(file, " copy.", ".", 1)
			// ommit the original photo
			submitted[file] = true
			skip[real] = true
		}
	}

	for _, file := range f.Photos {
		if skip[file] {
			continue
		}
		fi, err := os.Stat(filepath.Join(f.BasePath, file))
		if err != nil {
			log.Print(err)
			continue
		}
		p := Photo{
			Name:      file,
			Filename:  filepath.Join(f.BasePath, file),
			Size:      fi.Size(),
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

	sort.Slice(f.PhotoDetails, func(i, j int) bool { return f.PhotoDetails[i].Created.Before(f.PhotoDetails[j].Created) })

	return
}

func (f *FullComplaint) HasGPSInfo() bool {
	if f.Long != 0 && f.Lat != 0 {
		return true
	}
	f.ParsePhotos()
	for _, p := range f.PhotoDetails {
		if p.Lat != 0 && p.Long != 0 {
			return true
		}
	}
	return false
}
