package testhelpers

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

func BuildGoamzS3(accessKey, secretKey, endpoint string) *s3.S3 {
	auth := aws.Auth{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}

	region := aws.Region{
		Name:                 "aRegion",
		S3Endpoint:           endpoint,
		S3LocationConstraint: true,
	}

	return s3.New(auth, region)
}
