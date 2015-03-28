package commandline_test

import (
	"github.com/mitchellh/goamz/s3/s3test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestCommandline(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commandline Suite")
}

var s3Server *s3test.Server
var s3EndpointURL string
var cli string

var _ = BeforeSuite(func() {
	var err error
	s3Server, err = s3test.NewServer(nil)
	if err != nil {
		Expect(err).ToNot(HaveOccurred())
	}

	cli = buildCli()
	s3EndpointURL = s3Server.URL()
})

var _ = AfterSuite(func() {
	s3Server.Quit()
})

func buildCli() string {
	cliPath, err := gexec.Build("github.com/tscolari/s3up")
	Expect(err).ToNot(HaveOccurred())
	return cliPath
}
