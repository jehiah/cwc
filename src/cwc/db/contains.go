package db

import (
	"bytes"
	"io/ioutil"
	"strings"
)

func (d DB) ComplaintContains(c Complaint, pattern string) (bool, error) {
	f, err := d.Open(c)
	if err != nil {
		return false, err
	}
	defer f.Close()
	body, err := ioutil.ReadAll(f)
	if err != nil {
		return false, err
	}
	if bytes.Contains(body, []byte(pattern)) {
		return true, nil
	}
	return false, nil
}

func (f FullComplaint) Contains(pattern string) bool {
	return strings.Contains(f.Body, pattern)
}
