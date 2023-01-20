package db

import (
	"bytes"
	"io/ioutil"

	"github.com/jehiah/cwc/internal/complaint"
)

func (d LocalFilesystem) ComplaintContains(c complaint.Complaint, pattern string) (bool, error) {
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
