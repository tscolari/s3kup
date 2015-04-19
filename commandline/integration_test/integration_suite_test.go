package integration_test

import (
	"bytes"
	"os/exec"

	"github.com/mitchellh/goamz/s3"
	"github.com/mitchellh/goamz/s3/s3test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/tscolari/s3kup/s3/testhelpers"

	"testing"
)

func TestCommandline(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commandline > Integration Suite")
}

var s3Client *s3.S3
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
	cliPath, err := gexec.Build("github.com/tscolari/s3kup")
	Expect(err).ToNot(HaveOccurred())
	return cliPath
}

func runPipedCmdsAndReturnLastOutput(firstCmd, lastCmd *exec.Cmd) (string, error) {
	var outputBuffer bytes.Buffer
	var err error

	lastCmd.Stdin, _ = firstCmd.StdoutPipe()
	lastCmd.Stdout = &outputBuffer
	lastCmd.Stderr = &outputBuffer

	firstCmd.Start()
	lastCmd.Start()
	err = lastCmd.Wait()

	return outputBuffer.String(), err
}

func s3Bucket(accessKey, secretKey, bucketName string) *s3.Bucket {
	s3Client = testhelpers.BuildGoamzS3(accessKey, secretKey, s3Server.URL())
	return s3Client.Bucket(bucketName)
}
