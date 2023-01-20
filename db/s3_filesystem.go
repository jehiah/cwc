package db

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path"

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
		Prefix: prefix,
	}
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
	return nil, fmt.Errorf("not implemented")
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
	obj, err := d.client.GetObject(ctx, &s3.GetObjectInput{Bucket: &d.Bucket, Key: &p})
	if err != nil {
		return complaint.RawComplaint{}, err
	}
	defer obj.Body.Close()
	rc := complaint.RawComplaint{
		Complaint: c,
	}
	rc.Body, err = ioutil.ReadAll(obj.Body)
	return rc, err
}
func (d *S3DB) Attachments(c complaint.Complaint) ([]fs.DirEntry, error) {
	p := path.Join(d.Prefix, string(c))
	paginator := s3.NewListObjectsV2Paginator(d.client, &s3.ListObjectsV2Input{
		Bucket: &d.Bucket,
		Prefix: &p,
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, obj := range page.Contents {
			log.Printf("obj %s", *obj.Key)
		}
	}

	return nil, fmt.Errorf("not implemented")
}
func (d *S3DB) OpenAttachment(c complaint.Complaint, filename string) (fs.File, error) {
	return nil, fmt.Errorf("not implemented")
}
