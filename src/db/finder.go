package db

import (
	"bytes"
	"io/ioutil"
	"os"
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
	return o, nil
}
