package db

import (
	"bytes"
	"io/ioutil"
)

func (d DB) ComplaintContains(c Complaint, pattern string) (bool, error) {
	f, err := d.Open(c)
	if err != nil {
		return false, err
	}
	body, err := ioutil.ReadAll(f)
	if err != nil {
		return false, err
	}
	if bytes.Contains(body, []byte(pattern)) {
		return true, nil
	}
	return false, nil
}
