package db

import (
	"errors"
	"io"
	"io/fs"
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

func (d LocalFilesystem) Create(c complaint.Complaint) (*os.File, error) {
	err := os.MkdirAll(d.FullPath(c), os.ModePerm)
	if err != nil {
		return nil, err
	}
	fullPath := filepath.Join(d.FullPath(c), "notes.txt")
	return os.Create(fullPath)
}

func (d LocalFilesystem) Open(c complaint.Complaint) (*os.File, error) {
	fullPath := filepath.Join(d.FullPath(c), "notes.txt")
	return os.Open(fullPath)
}

func (d LocalFilesystem) Append(c complaint.Complaint, s string) error {
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
func (d LocalFilesystem) Exists(c complaint.Complaint) (bool, error) {
	fullPath := filepath.Join(d.FullPath(c), "notes.txt")
	_, err := os.Stat(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	} else if err != nil {
		return false, nil
	}
	return true, nil
}

func (d LocalFilesystem) Read(c complaint.Complaint) (complaint.RawComplaint, error) {
	f, err := d.Open(c)
	if err != nil {
		return complaint.RawComplaint{}, err
	}
	defer f.Close()
	body, err := ioutil.ReadAll(f)
	if err != nil {
		return complaint.RawComplaint{}, err
	}
	return complaint.RawComplaint{
		Complaint: c,
		Body:      body,
	}, nil

}
func (d LocalFilesystem) Attachments(c complaint.Complaint) ([]fs.DirEntry, error) {
	files, err := fs.ReadDir(os.DirFS(string(d)), string(c))
	if err != nil {
		return nil, err
	}
	var out []fs.DirEntry
	for _, f := range files {
		switch f.Name() {
		case "notes.txt", ".DS_Store":
			continue
		}
		out = append(out, f)
	}
	return out, nil
}

func (d LocalFilesystem) OpenAttachment(c complaint.Complaint, filename string) (fs.File, error) {
	fullPath := filepath.Join(d.FullPath(c), filename)
	return os.Open(fullPath)
}

func (d LocalFilesystem) FullComplaint(c complaint.Complaint) (*complaint.FullComplaint, error) {
	rc, err := d.Read(c)
	if err != nil {
		return nil, err
	}

	files, err := d.Attachments(c)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, f := range files {
		filenames = append(filenames, f.Name())
	}
	return complaint.ParseComplaint(rc, filenames)
}
