package db

import (
	"io"
	"io/fs"
	"os"

	"github.com/jehiah/cwc/internal/complaint"
)

type ReadOnly interface {
	Index() ([]complaint.Complaint, error)
	Exists(c complaint.Complaint) (bool, error)
	FullPath(c complaint.Complaint) string
	Latest() (complaint.Complaint, error)

	Find(pattern string) ([]complaint.Complaint, error)
	FullComplaint(c complaint.Complaint) (*complaint.FullComplaint, error)
	ComplaintContains(c complaint.Complaint, pattern string) (bool, error)

	Read(complaint.Complaint) (complaint.RawComplaint, error)
	Attachments(complaint.Complaint) ([]fs.DirEntry, error)
	OpenAttachment(c complaint.Complaint, filename string) (io.ReadCloser, error)
}

type Write interface {
	Append(c complaint.Complaint, s string) error
	Create(c complaint.Complaint) (*os.File, error)
	Open(c complaint.Complaint) (*os.File, error)
}

type ReadWrite interface {
	ReadOnly
	Write
}

type Interactive interface {
	ShowInFinder(c complaint.Complaint) error
	ShowInEditor(c complaint.Complaint) error
}
