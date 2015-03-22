package integration_test

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3/s3test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	goamzs3 "github.com/mitchellh/goamz/s3"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var s3Server *s3test.Server
var s3EndpointURL string

var _ = BeforeSuite(func() {
	var err error
	config := &s3test.Config{}
	s3Server, err = s3test.NewServer(config)
	if err != nil {
		Expect(err).ToNot(HaveOccurred())
	}

	s3EndpointURL = s3Server.URL()
})

var _ = AfterSuite(func() {
	s3Server.Quit()
})

func buildGoamzS3(accessKey, secretKey, regionName string) *goamzs3.S3 {
	auth := aws.Auth{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}

	region := aws.Region{
		Name:                 regionName,
		S3Endpoint:           s3EndpointURL,
		S3LocationConstraint: true,
	}

	return goamzs3.New(auth, region)
}
