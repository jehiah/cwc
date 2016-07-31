package db

import (
	"bytes"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"sort"
)

// DB is the directory path containing a cwc repo
type DB string

var Default DB

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err.Error())
	}
	Default = DB(filepath.Join(usr.HomeDir, "Documents", "cyclists_with_cameras"))
}

// Find finds the CWC reports that have a given pattern in them
func (d DB) Find(pattern string) ([]Complaint, error) {
	matches, err := filepath.Glob(filepath.Join(string(d), "*", "notes.txt"))
	if err != nil {
		return nil, err
	}
	var found []Complaint
	for _, f := range matches {
		body, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}
		if bytes.Contains(body, []byte(pattern)) {
			found = append(found, Complaint(filepath.Base(filepath.Dir(f))))
		}
	}
	sort.Sort(sort.Reverse(complaintsByAge(found)))
	return found, nil
}
