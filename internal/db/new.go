package db

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jehiah/cwc/internal/complaint"
)

// New constructs a new complaint directory
func (d DB) New(dt time.Time, license string) (complaint.Complaint, error) {
	c := fmt.Sprintf("%s_%s", dt.Format("20060102_1504"), license)
	fullPath := filepath.Join(string(d), c)
	err := os.MkdirAll(fullPath, os.ModePerm)
	return complaint.Complaint(c), err
}

// FullPath reuturns the absolute path to the complaint directory
func (d DB) FullPath(c complaint.Complaint) string {
	return filepath.Join(string(d), string(c))
}
