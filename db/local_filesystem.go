package db

import (
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"sort"

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

func (d LocalFilesystem) FullComplaint(c complaint.Complaint) (*complaint.FullComplaint, error) {
	return FullComplaint(d, c)
}
func (d LocalFilesystem) ComplaintContains(c complaint.Complaint, pattern string) (bool, error) {
	return ComplaintContains(d, c, pattern)
}
func (d LocalFilesystem) Find(pattern string) ([]complaint.Complaint, error) { return Find(d, pattern) }

// FullPath reuturns the absolute path to the complaint directory
func (d LocalFilesystem) FullPath(c complaint.Complaint) string {
	return filepath.Join(string(d), string(c))
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

// Index returns all complaints in a DB directory. It assumes all top level directories are complaints
func (d LocalFilesystem) Index() ([]complaint.Complaint, error) {
	var o []complaint.Complaint
	f, err := os.Open(string(d))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	finfos, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	for _, fi := range finfos {
		if !fi.IsDir() {
			continue
		}
		o = append(o, complaint.Complaint(fi.Name()))
	}
	sort.Sort(sort.Reverse(complaint.ComplaintsByAge(o)))
	return o, nil
}

// Latest - TODO: rename LastModified?
func (d LocalFilesystem) Latest() (complaint.Complaint, error) {
	var latest os.FileInfo
	f, err := os.Open(string(d))
	if err != nil {
		return complaint.Complaint(""), err
	}
	defer f.Close()
	finfos, err := f.Readdir(-1)
	if err != nil {
		return complaint.Complaint(""), err
	}
	for _, fi := range finfos {
		if !fi.IsDir() {
			continue
		}
		if latest == nil || fi.ModTime().After(latest.ModTime()) {
			latest = fi
		}
	}
	if latest == nil {
		return complaint.Complaint(""), nil
	}
	return complaint.Complaint(latest.Name()), nil
}
