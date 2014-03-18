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
	"strings"
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
	S3ACLPrivate 				string = "private"
	S3ACLPublicRead 			string = "public-read"
	S3ACLPublicReadWrite 		string = "public-read-write"
	S3ACLAuthenticatedRead 		string = "authenticated-read"
	S3ACLBucketOwnerRead 		string = "bucket-owner-read"
	S3ACLBucketOwnerFullControl string = "bucket-owner-full-control"
)

func InitS3(n int) {
	flag.IntVar(&S3Workers, "s3workers", n, "Number of S3 uploaders.")
}

func StartS3Uploaders(
	work chan interface{}, errorHander func (S3Upload, error),
) {
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
				fmt.Println("worker waiting", id)
				upload := (<- work).(S3Upload)
				fmt.Println("worker got upload", id, upload.Path)
				err := UploadToS3(upload, client, keys)
				if err != nil {
					errorHander(upload, err)
				}
			}
		}(i)
		s3wg.Add(1)
	}
}

func UploadToS3(upload S3Upload, client http.Client, keys s3.Keys) error {
	url := fmt.Sprintf(
		"https://%s.s3.amazonaws.com%s", upload.Bucket, upload.Path,
	)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(upload.Content))
	if err != nil {
		fmt.Println("err1", err, url)
		return err
	}

	req.ContentLength = int64(len(upload.Content))
	req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	if upload.ACL == "" {
		req.Header.Set("X-Amz-Acl", S3ACLPublicRead)
	} else {
		req.Header.Set("X-Amz-Acl", upload.ACL)
	}
	if upload.MimeType != "" {
		req.Header.Set("Content-Type", upload.MimeType)
	}
	s3.Sign(req, keys)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("err2", err, url)
		res, err = http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("err2.2", err, url)
			filepath := fmt.Sprintf("/tmp/%s", strings.Split(upload.Path, "/")[2])
			fmt.Println("create file", filepath)
			file, err := os.Create(filepath)
			if err != nil {
				fmt.Println("cant open file", filepath)
			}
			file.Write(upload.Content)
			file.Close()
			return err
		}
	}

	if res.StatusCode > 399 {
		fmt.Println("err3", err, url)
		return errors.New(
			fmt.Sprintf("Upload Failed: %d for %s.", res.StatusCode, url),
		)
	}

	_, err = io.Copy(ioutil.Discard, res.Body)
	if err != nil {
		fmt.Println("err5", err, url)
		return err
	}

	err = res.Body.Close()
	if err != nil {
		fmt.Println("err5", err, url)
		return err
	}

	return nil
}

func WaitForS3Uploaders() {
	s3wg.Wait()
}