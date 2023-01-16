package db

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jehiah/cwc/internal/complaint"
)

var editors []string = []string{
	"mate",
	"atom",
}

func (d DB) Edit(c complaint.Complaint) error {
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

func (d DB) Create(c complaint.Complaint) (*os.File, error) {
	fullPath := filepath.Join(d.FullPath(c), "notes.txt")
	return os.Create(fullPath)
}

func (d DB) Open(c complaint.Complaint) (*os.File, error) {
	fullPath := filepath.Join(d.FullPath(c), "notes.txt")
	return os.Open(fullPath)
}

func (d DB) ShowInFinder(c complaint.Complaint) error {
	return exec.Command("/usr/bin/open", d.FullPath(c)).Run()
}

func (d DB) Append(c complaint.Complaint, s string) error {
	fullPath := filepath.Join(d.FullPath(c), "notes.txt")
	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY, 066)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.WriteString(f, s)
	return err
}

// Exists checks if the complaint exists in DB
func (d DB) Exists(c complaint.Complaint) (bool, error) {
	fullPath := filepath.Join(d.FullPath(c), "notes.txt")
	_, err := os.Stat(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	} else if err != nil {
		return false, nil
	}
	return true, nil
}

func (d DB) FullComplaint(c complaint.Complaint) (*complaint.FullComplaint, error) {
	f, err := d.Open(c)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	body, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	path := d.FullPath(c)
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	files, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	return complaint.ParseComplaint(c, body, path, files)
}
