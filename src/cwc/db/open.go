package db

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

var editors []string = []string{
	"/usr/local/bin/mate",
	"/usr/local/bin/atom",
	"/Applications/Sublime Text.app/Contents/SharedSupport/bin/subl",
	"/Applications/TextEdit.app",
}

func (d DB) Edit(c Complaint) error {
	for _, editor := range editors {
		if _, err := os.Stat(editor); err != nil {
			continue
		}
		fullPath := filepath.Join(d.FullPath(c), "notes.txt")
		return exec.Command(editor, fullPath).Run()
	}
	return errors.New("No editor found")
}

func (d DB) Create(c Complaint) (*os.File, error) {
	fullPath := filepath.Join(d.FullPath(c), "notes.txt")
	return os.Create(fullPath)
}

func (d DB) Open(c Complaint) (*os.File, error) {
	fullPath := filepath.Join(d.FullPath(c), "notes.txt")
	return os.Open(fullPath)
}

func (d DB) ShowInFinder(c Complaint) error {
	return exec.Command("/usr/bin/open", d.FullPath(c)).Run()
}
