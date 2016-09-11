package db

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
)

func (d DB) ComplaintContains(c Complaint, pattern string) (bool, error) {
	f := filepath.Join(string(d), string(c), "notes.txt")
	body, err := ioutil.ReadFile(f)
	if err != nil {
		return false, err
	}
	if bytes.Contains(body, []byte(pattern)) {
		return true, nil
	}
	return false, nil
}
