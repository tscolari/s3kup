package integration_test

import (
	"github.com/mitchellh/goamz/s3/s3test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
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
