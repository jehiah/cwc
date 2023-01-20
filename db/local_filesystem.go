package db

import (
	"os/user"
	"path/filepath"

	"github.com/jehiah/cwc/internal/complaint"
)

// DB is the directory path containing a cwc repo
type LocalFilesystem string

// Default is the default DB at ~/Documents/cyclists_with_cameras
var Default LocalFilesystem

func init() {
	defer func() {
		recover()
	}()

	usr, err := user.Current()
	if err != nil {
		panic(err.Error())
	}
	Default = LocalFilesystem(filepath.Join(usr.HomeDir, "Documents", "cyclists_with_cameras"))
}

// FullPath reuturns the absolute path to the complaint directory
func (d LocalFilesystem) FullPath(c complaint.Complaint) string {
	return filepath.Join(string(d), string(c))
}
