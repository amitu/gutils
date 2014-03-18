package gutils

import (
	"io"
	"os"
	"fmt"
	"flag"
	"sync"
	"time"
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"github.com/kr/s3"
)

type S3Upload struct {
	// if passed this will be uploaded
	Content []byte
	// either content of file name must be passed
	FileName string
	// s3 bucket, can be left empty if default bucket is set
	Bucket string
	MimeType string
	// path within the bucket
	Path string
	ACL string
}

var (
	S3Workers int = 2
	s3wg sync.WaitGroup
)

const (
	S3ACLPrivate string = "private"
	 // | public-read | public-read-write | authenticated-read | bucket-owner-read | bucket-owner-full-control
)

func InitS3(n int) {
	flag.IntVar(&S3Workers, "s3workers", n, "Number of S3 uploaders.")
}

func StartS3Uploaders(work chan interface{}, errorHander func (err error)) {
	client := http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: S3Workers,
		},
	}

	keys := s3.Keys{
	    AccessKey: os.Getenv("S3_ACCESS_KEY"),
	    SecretKey: os.Getenv("S3_SECRET_KEY"),
	}

	for i := 0; i < S3Workers; i++ {
		go func(id int) {
			for {
				upload := (<- work).(S3Upload)
				err := UploadToS3(upload, client, keys)
				if err != nil {
					errorHander(err)
				}
			}
		}(i)
		s3wg.Add(1)
	}
}

func UploadToS3(upload S3Upload, client http.Client, keys s3.Keys) error {
	url := fmt.Sprintf("https://s3.amazonaws.com/%s", upload.Bucket)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(upload.Content))
	if err != nil {
		return err
	}

	req.ContentLength = int64(len(upload.Content))
	req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	req.Header.Set("X-Amz-Acl", upload.ACL)
	s3.Sign(req, keys)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode > 399 {
		return errors.New(
			fmt.Sprintf("Upload Failed: %d for %s.", res.StatusCode, url),
		)
	}

	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()
	return nil
}

func WaitForS3Uploaders() {
	s3wg.Wait()
}