package gutils

import (
	"os"
	"io"
	"log"
	"fmt"
	"flag"
	"sync"
	"time"
	"bytes"
	"errors"
	"net/http"
	"io/ioutil"
	"github.com/kr/s3"
)

type S3UploadContext interface {
	OnError(S3Upload, error)
}

type S3Upload struct {
	// if passed this will be uploaded
	Content  []byte
	// either content of file name must be passed
	FileName string
	// s3 bucket, can be left empty if default bucket is set
	Bucket   string
	MimeType string
	// path within the bucket
	Path     string
	ACL      string
	Context  S3UploadContext
}

var (
	S3Workers int = 2
	s3wg      sync.WaitGroup
	s3keys    s3.Keys
)

const (
	S3ACLPrivate 				string = "private"
	S3ACLPublicRead 			string = "public-read"
	S3ACLPublicReadWrite 		string = "public-read-write"
	S3ACLAuthenticatedRead 		string = "authenticated-read"
	S3ACLBucketOwnerRead 		string = "bucket-owner-read"
	S3ACLBucketOwnerFullControl string = "bucket-owner-full-control"
)

func InitS3 (n int) {
	// TODO: probably command line flag should not be set by a library like
	// this and should be relied on application to set number of workers based
	// on whatever, eg config file
	flag.IntVar(&S3Workers, "s3workers", n, "Number of S3 uploaders.")
	s3keys = s3.Keys{
	    AccessKey: os.Getenv("S3_ACCESS_KEY"),
	    SecretKey: os.Getenv("S3_SECRET_KEY"),
	}
}

func InitS3WithKeys (n int, access, secret string) {
	InitS3(n)
	s3keys = s3.Keys{
	    AccessKey: access,
	    SecretKey: secret,
	}
}

func StartS3Uploaders (work chan interface{}) {

	client := http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: S3Workers,
		},
	}

	for i := 0; i < S3Workers; i++ {
		go func() {
			for {
				job, found := <- work
				if !found {
					// Channel closed
					s3wg.Done()
					return
				}
				upload, ok := job.(S3Upload)
				if !ok {
					// TODO: this is major, how to properly handle it?
					fmt.Println("gutils.S3: [CRITICAL] Bad job found, dropped.")
					continue
				}
				err := UploadToS3(upload, client)
				if err != nil {
					upload.Context.OnError(upload, err)
				}
			}
		}()
		s3wg.Add(1)
	}
}

func UploadToS3(upload S3Upload, client http.Client) error {
	url := fmt.Sprintf(
		"https://%s.s3.amazonaws.com%s", upload.Bucket, upload.Path,
	)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(upload.Content))
	if err != nil {
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
	s3.Sign(req, s3keys)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode > 399 {
		return errors.New(
			fmt.Sprintf("Upload Failed: %d for %s.", res.StatusCode, url),
		)
	}

	_, err = io.Copy(ioutil.Discard, res.Body)
	if err != nil {
		return err
	}

	err = res.Body.Close()
	if err != nil {
		return err
	}

	log.Println("gutils.S3: Uploaded " + url)
	return nil
}

func WaitForS3Uploaders() {
	s3wg.Wait()
}