package db

import (
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jehiah/cwc/exif"
	"github.com/jehiah/cwc/internal/complaint"
)

func LoadPhotos(d ReadOnly, f *complaint.FullComplaint) ([]complaint.Photo, error) {
	// if len(f.PhotoDetails) > 0 && len(f.Photos) > 0 {
	// 	// idempotent
	// 	return nil, nil
	// }

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

	var details []complaint.Photo
	for _, file := range f.Photos {
		if skip[file] {
			continue
		}

		f, err := d.OpenAttachment(f.Complaint, file)
		if err != nil {
			return nil, err
		}
		// TODO: get photo size
		// fi, err := f.Stat()
		// if err != nil {
		// 	return nil, err
		// }

		p := complaint.Photo{
			Name: file,
			// Size:      fi.Size(),
			Submitted: submitted[file],
		}
		x, err := exif.Parse(f)
		if err != nil {
			// TODO: use file created timestamp
			log.Printf("%s", err)
		} else {
			p.Created = x.Created
			p.Lat, p.Long = x.Lat, x.Long
		}
		details = append(details, p)
	}

	sort.Slice(details, func(i, j int) bool { return details[i].Created.Before(details[j].Created) })

	return details, nil
}
