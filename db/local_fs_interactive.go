package db

import (
	"errors"
	"os/exec"
	"path/filepath"

	"github.com/jehiah/cwc/internal/complaint"
)

var editors []string = []string{
	"mate",
	"atom",
}

func (d LocalFilesystem) ShowInEditor(c complaint.Complaint) error {
	for _, editor := range editors {
		if e, err := exec.LookPath(editor); err != nil {
			continue
		} else {
			editor = e
		}
		fullPath := filepath.Join(d.FullPath(c), "notes.txt")
		return exec.Command(editor, fullPath).Run()
	}
	return errors.New("No editor found")
}

func (d LocalFilesystem) ShowInFinder(c complaint.Complaint) error {
	return exec.Command("/usr/bin/open", d.FullPath(c)).Run()
}
