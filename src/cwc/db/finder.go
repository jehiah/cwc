package db

import (
	"os"
	"path/filepath"
	"sort"
)

// Find finds the CWC reports that have a given pattern in them
func (d DB) Find(pattern string) ([]Complaint, error) {
	matches, err := filepath.Glob(filepath.Join(string(d), "*", "notes.txt"))
	if err != nil {
		return nil, err
	}
	var found []Complaint
	for _, f := range matches {
		complaint := Complaint(filepath.Base(filepath.Dir(f)))
		if ok, err := d.ComplaintContains(complaint, pattern); err != nil {
			return nil, err
		} else if ok {
			found = append(found, complaint)
		}
	}
	sort.Sort(sort.Reverse(complaintsByAge(found)))
	return found, nil
}

// All returns all complaints in a DB directory. It assumes all top level directories are complaints
func (d DB) All() ([]Complaint, error) {
	var o []Complaint
	f, err := os.Open(string(d))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	finfos, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	for _, fi := range finfos {
		if !fi.IsDir() {
			continue
		}
		o = append(o, Complaint(fi.Name()))
	}
	sort.Sort(sort.Reverse(complaintsByAge(o)))
	return o, nil
}
