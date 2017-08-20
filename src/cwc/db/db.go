package db

import (
	"os/user"
	"path/filepath"
)

// DB is the directory path containing a cwc repo
type DB string

// Default is the default DB at ~/Documents/cyclists_with_cameras
var Default DB

func init() {
	// in Go 1.8 this isn't implemented in linux
	defer func() {
		recover()
	}()

	usr, err := user.Current()
	if err != nil {
		panic(err.Error())
	}
	Default = DB(filepath.Join(usr.HomeDir, "Documents", "cyclists_with_cameras"))
}
