package db

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// New constructs a new complaint directory
func (d DB) New(dt time.Time, license string) (Complaint, error) {
	complaint := fmt.Sprintf("%s_%s", dt.Format("20060102_1504"), license)
	fullPath := filepath.Join(string(d), complaint)
	err := os.MkdirAll(fullPath, os.ModePerm)
	return Complaint(complaint), err
}

// FullPath reuturns the absolute path to the complaint directory
func (d DB) FullPath(c Complaint) string {
	return filepath.Join(string(d), string(c))
}
