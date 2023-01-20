package db

import (
	"bytes"

	"github.com/jehiah/cwc/internal/complaint"
)

func FullComplaint(d ReadOnly, c complaint.Complaint) (*complaint.FullComplaint, error) {
	rc, err := d.Read(c)
	if err != nil {
		return nil, err
	}

	files, err := d.Attachments(c)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, f := range files {
		filenames = append(filenames, f.Name())
	}
	return complaint.ParseComplaint(rc, filenames)
}

// Find finds the CWC reports that have a given pattern in them
func Find(d ReadOnly, pattern string) ([]complaint.Complaint, error) {
	all, err := d.Index()
	if err != nil {
		return nil, err
	}
	var found []complaint.Complaint
	for _, c := range all {
		if ok, err := d.ComplaintContains(c, pattern); err != nil {
			return nil, err
		} else if ok {
			found = append(found, c)
		}

	}
	return found, nil
}

func ComplaintContains(d ReadOnly, c complaint.Complaint, pattern string) (bool, error) {
	rc, err := d.Read(c)
	if err != nil {
		return false, err
	}
	return bytes.Contains(rc.Body, []byte(pattern)), nil
}
