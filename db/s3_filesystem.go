package db

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jehiah/cwc/internal/complaint"
)

type S3DB struct {
	client *s3.Client
	Bucket string
	Prefix string
}

func NewS3DB(client *s3.Client, bucket, prefix string) *S3DB {
	return &S3DB{
		client: client,
		Bucket: bucket,
		Prefix: strings.TrimLeft(prefix, "/"),
	}
}

func (d *S3DB) Append(c complaint.Complaint, s string) error { return fmt.Errorf("not implemented") }
func (d *S3DB) Create(c complaint.Complaint) (*os.File, error) {
	return nil, fmt.Errorf("not implemented")
}
func (d *S3DB) Open(c complaint.Complaint) (*os.File, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d *S3DB) Find(pattern string) ([]complaint.Complaint, error) {
	return Find(d, pattern)
}
func (d *S3DB) FullComplaint(c complaint.Complaint) (*complaint.FullComplaint, error) {
	return FullComplaint(d, c)
}
func (d *S3DB) ComplaintContains(c complaint.Complaint, pattern string) (bool, error) {
	return ComplaintContains(d, c, pattern)
}

func (d *S3DB) Index() ([]complaint.Complaint, error) {
	fs, err := d.Attachments(complaint.Complaint(""))
	if err != nil {
		return nil, err
	}
	var out []complaint.Complaint
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), "/notes.txt") {
			n := strings.Split(strings.TrimSuffix(f.Name(), "/notes.txt"), "/")
			c := complaint.Complaint(n[len(n)-1])
			log.Printf("got %q", string(c))
			out = append(out, c)
		}
	}
	return out, nil
}

func (d *S3DB) Exists(c complaint.Complaint) (bool, error) {
	ctx := context.Background()
	p := path.Join(d.Prefix, string(c), "notes.txt")
	_, err := d.client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: &d.Bucket, Key: &p})
	if err != nil {
		var re *awshttp.ResponseError
		if errors.As(err, &re) && re.Response.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func (d *S3DB) FullPath(c complaint.Complaint) string {
	return ""
}
func (d *S3DB) Latest() (complaint.Complaint, error) {
	return "", fmt.Errorf("not implemented")
}
func (d *S3DB) Read(c complaint.Complaint) (complaint.RawComplaint, error) {
	ctx := context.Background()
	p := path.Join(d.Prefix, string(c), "notes.txt")
	log.Printf("Read %q", p)
	obj, err := d.client.GetObject(ctx, &s3.GetObjectInput{Bucket: &d.Bucket, Key: &p})
	if err != nil {
		return complaint.RawComplaint{}, err
	}
	defer obj.Body.Close()
	rc := complaint.RawComplaint{
		Complaint:    c,
		LastModified: *obj.LastModified,
	}
	rc.Body, err = ioutil.ReadAll(obj.Body)
	return rc, err
}

type s3Info struct {
	name    string
	size    int64
	modTime time.Time
}

func (s s3Info) Name() string {
	return s.name
}
func (s s3Info) Size() int64 {
	return s.size
}
func (s s3Info) Mode() fs.FileMode {
	return fs.ModePerm
}
func (s s3Info) ModTime() time.Time {
	return s.modTime
}
func (s s3Info) IsDir() bool {
	return false
}
func (s s3Info) Sys() any {
	return nil
}

type s3Obj struct {
	name    string
	size    int64
	modTime time.Time
}

func (s s3Obj) Name() string      { return s.name }
func (s s3Obj) IsDir() bool       { return false }
func (s s3Obj) Type() fs.FileMode { return fs.ModePerm }
func (s s3Obj) Info() (fs.FileInfo, error) {
	return s3Info{
		name:    s.name,
		size:    s.size,
		modTime: s.modTime,
	}, nil
}

func (d *S3DB) Attachments(c complaint.Complaint) ([]fs.DirEntry, error) {
	ctx := context.Background()
	p := path.Join(d.Prefix, string(c))
	// log.Printf("list prefix %q", p)
	paginator := s3.NewListObjectsV2Paginator(d.client, &s3.ListObjectsV2Input{
		Bucket: &d.Bucket,
		Prefix: &p,
	})

	var items []fs.DirEntry
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, obj := range page.Contents {
			log.Printf("got obj %s", *obj.Key)
			items = append(items, s3Obj{
				name:    strings.TrimPrefix(*obj.Key, p),
				size:    *obj.Size,
				modTime: *obj.LastModified,
			})
		}
	}
	return items, nil
}
func (d *S3DB) OpenAttachment(c complaint.Complaint, filename string) (io.ReadCloser, error) {
	ctx := context.Background()
	p := path.Join(d.Prefix, string(c), filename)
	log.Printf("OpenAttachment %q", p)
	obj, err := d.client.GetObject(ctx, &s3.GetObjectInput{Bucket: &d.Bucket, Key: &p})
	if err != nil {
		return nil, err
	}
	return obj.Body, nil
}
