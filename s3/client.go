package s3

import (
	"github.com/mitchellh/goamz/aws"
	goamzs3 "github.com/mitchellh/goamz/s3"
)

type Client struct {
	s3     *goamzs3.S3
	bucket *goamzs3.Bucket
}

func New(accessKey, secretKey, bucketName, endPointURL string) *Client {
	auth := aws.Auth{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}

	region := getRegion(endPointURL)

	s3 := goamzs3.New(auth, region)
	bucket := s3.Bucket(bucketName)

	return &Client{
		s3:     s3,
		bucket: bucket,
	}
}

func getRegion(endPointURL string) aws.Region {
	for _, region := range aws.Regions {
		if region.S3Endpoint == endPointURL {
			return region
		}
	}

	return aws.Region{
		Name:                 "custom",
		S3Endpoint:           endPointURL,
		S3LocationConstraint: true,
	}
}

func (c *Client) Store(path string, fileContent []byte) error {
	return c.bucket.Put(path, fileContent, "", "")
}

func (c *Client) Delete(path string) error {
	return c.bucket.Del(path)
}

func (c *Client) List(path string) (Versions, error) {
	resp, err := c.bucket.List(path+"/", "", "", 100)
	if err != nil {
		return []Version{}, err
	}

	files := []Version{}
	for _, file := range resp.Contents {
		version, err := NewVersion(file)
		if err != nil {
			return Versions{}, err
		}
		files = append(files, version)
	}

	return files, nil
}

func (c *Client) Get(path string) ([]byte, error) {
	return c.bucket.Get(path)
}
