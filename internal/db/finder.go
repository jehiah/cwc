package db

import (
	"os"
	"sort"

	"github.com/jehiah/cwc/internal/complaint"
)

// Find finds the CWC reports that have a given pattern in them
func (d LocalFilesystem) Find(pattern string) ([]complaint.Complaint, error) {
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

// Index returns all complaints in a DB directory. It assumes all top level directories are complaints
func (d LocalFilesystem) Index() ([]complaint.Complaint, error) {
	var o []complaint.Complaint
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
		o = append(o, complaint.Complaint(fi.Name()))
	}
	sort.Sort(sort.Reverse(complaint.ComplaintsByAge(o)))
	return o, nil
}

func (d LocalFilesystem) Latest() (complaint.Complaint, error) {
	var latest os.FileInfo
	f, err := os.Open(string(d))
	if err != nil {
		return complaint.Complaint(""), err
	}
	defer f.Close()
	finfos, err := f.Readdir(-1)
	if err != nil {
		return complaint.Complaint(""), err
	}
	for _, fi := range finfos {
		if !fi.IsDir() {
			continue
		}
		if latest == nil || fi.ModTime().After(latest.ModTime()) {
			latest = fi
		}
	}
	if latest == nil {
		return complaint.Complaint(""), nil
	}
	return complaint.Complaint(latest.Name()), nil
}
