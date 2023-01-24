package db

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/jehiah/cwc/internal/complaint"
)

type S3Cache struct {
	db         *S3DB
	complaints map[complaint.Complaint]*complaintAttachments
}

type complaintAttachments struct {
	complaint.RawComplaint
	attachments []s3Obj
}

func NewS3Cache(db *S3DB) *S3Cache {

	c := &S3Cache{
		db:         db,
		complaints: make(map[complaint.Complaint]*complaintAttachments),
	}
	err := c.Load()
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func (cache *S3Cache) Load() error {
	files, err := cache.db.Attachments("")
	if err != nil {
		return err
	}

	// build all the attachments
	for _, f := range files {
		if strings.Count(f.Name(), "/") != 1 {
			continue
		}
		c := complaint.Complaint(strings.Split(f.Name(), "/")[0])
		fn := strings.Split(f.Name(), "/")[1]

		ca, ok := cache.complaints[c]
		if !ok {
			ca = &complaintAttachments{}
			cache.complaints[c] = ca
		}
		switch fn {
		case "notes.txt":
			continue
		}
		s := f.(s3Obj)
		s.name = fn
		ca.attachments = append(ca.attachments, s)
	}
	var wg sync.WaitGroup
	wg.Add(len(cache.complaints))
	for c, ca := range cache.complaints {
		c, ca := c, ca
		go func() {
			defer wg.Done()
			var err error
			ca.RawComplaint, err = cache.db.Read(c)
			if err != nil {
				log.Printf("%s %#v", c, err)
			}
		}()
	}
	wg.Wait()
	return nil
}

func (cache *S3Cache) Index() ([]complaint.Complaint, error) {
	var out []complaint.Complaint
	for c := range cache.complaints {
		out = append(out, c)
	}
	return out, nil
}

func (cache *S3Cache) Exists(c complaint.Complaint) (bool, error) {
	_, ok := cache.complaints[c]
	return ok, nil
}

func (cache *S3Cache) Read(c complaint.Complaint) (complaint.RawComplaint, error) {
	e, ok := cache.complaints[c]
	if !ok {
		return complaint.RawComplaint{}, fmt.Errorf("not found %q", string(c))
	}
	return e.RawComplaint, nil
}

func (cache *S3Cache) Attachments(c complaint.Complaint) ([]fs.DirEntry, error) {
	e, ok := cache.complaints[c]
	if !ok {
		return nil, fmt.Errorf("not found %q", string(c))
	}
	var o []fs.DirEntry
	for _, f := range e.attachments {
		o = append(o, f)
	}
	return o, nil
}

func (cache *S3Cache) OpenAttachment(c complaint.Complaint, filename string) (io.ReadCloser, error) {
	_, ok := cache.complaints[c]
	if !ok {
		return nil, fmt.Errorf("not found %q", string(c))
	}
	// TODO: if filename is present
	return cache.db.OpenAttachment(c, filename)
}

func (d *S3Cache) FullPath(c complaint.Complaint) string {
	return ""
}
func (d *S3Cache) Latest() (complaint.Complaint, error) {
	return "", fmt.Errorf("not implemented")
}

func (d *S3Cache) Append(c complaint.Complaint, s string) error { return fmt.Errorf("not implemented") }
func (d *S3Cache) Create(c complaint.Complaint) (*os.File, error) {
	return nil, fmt.Errorf("not implemented")
}
func (d *S3Cache) Open(c complaint.Complaint) (*os.File, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d *S3Cache) Find(pattern string) ([]complaint.Complaint, error) {
	return Find(d, pattern)
}
func (d *S3Cache) FullComplaint(c complaint.Complaint) (*complaint.FullComplaint, error) {
	return FullComplaint(d, c)
}
func (d *S3Cache) ComplaintContains(c complaint.Complaint, pattern string) (bool, error) {
	return ComplaintContains(d, c, pattern)
}
