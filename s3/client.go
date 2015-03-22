package s3

import (
	"github.com/mitchellh/goamz/aws"
	goamzs3 "github.com/mitchellh/goamz/s3"
)

type Client struct {
	s3     *goamzs3.S3
	bucket *goamzs3.Bucket
}

func New(accessKeyID, accessKeySecret, bucketName, regionName, endPointURL string) *Client {
	auth := aws.Auth{
		AccessKey: accessKeyID,
		SecretKey: accessKeySecret,
	}

	region := aws.Region{
		Name:                 regionName,
		S3LocationConstraint: true,
		S3Endpoint:           endPointURL,
	}

	s3 := goamzs3.New(auth, region)
	bucket := s3.Bucket(bucketName)
	bucket.PutBucket(goamzs3.Private)

	return &Client{
		s3:     s3,
		bucket: bucket,
	}
}

func (c *Client) Store(path string, fileContent []byte) error {
	return c.bucket.Put(path, fileContent, "", "")
}

func (c *Client) Delete(path string) error {
	return c.bucket.Del(path)
}

func (c *Client) List(path string) ([]string, error) {
	resp, err := c.bucket.List(path, "", "", 100)
	if err != nil {
		return []string{}, err
	}

	files := []string{}
	for _, file := range resp.Contents[1:] {
		files = append(files, file.Key)
	}

	return files, nil
}
