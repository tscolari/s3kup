package commandline_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"os/exec"

	"github.com/mitchellh/goamz/s3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tscolari/s3up/s3/testhelpers"
)

var s3Client *s3.S3

var _ = Describe("Cli", func() {

	const (
		accessKey      string = "my_id"
		secretKey      string = "my_secret"
		regionName     string = "my_region"
		versionsToKeep int    = 3
		backupName     string = "my-backup"
	)

	var backupCmd *exec.Cmd
	var inputCmd *exec.Cmd
	var bucketName string

	BeforeEach(func() {
		bucketName = fmt.Sprintf("bucket%d", rand.Int())
		inputCmd = exec.Command("echo", "'store my data'")
	})

	Describe("storing backups", func() {
		BeforeEach(func() {
			backupCmd = exec.Command(cli, "-i", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName, "--no-ssl", "-k", "3")
		})

		Context("when the correct args are given", func() {
			It("exits with failure if no input is given through a pipe", func() {
				output, err := backupCmd.Output()
				Expect(output).To(MatchRegexp("not using pipeline"))
				Expect(err).To(HaveOccurred())
			})

			It("stores the backup with the correct name", func() {
				bucket := s3Bucket(accessKey, secretKey, bucketName)
				bucket.PutBucket("")

				_, err := runPipedCmdsAndReturnLastOutput(inputCmd, backupCmd)
				Expect(err).ToNot(HaveOccurred())

				resp, err := bucket.List(backupName, "", "", 100)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(resp.Contents)).To(Equal(1))
			})

			It("keeps only the number of versions specified", func() {
				exec.Command(cli, "-i", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName, "--no-ssl", "-k", "2").Run()

			})
		})

		Context("when there is invalid or missing args", func() {
			It("fails if no file name is given", func() {
			})

			It("fails if no bucket name is given", func() {

			})

			It("fails if no access key is given", func() {

			})

			It("fails if no secret key is given", func() {

			})

			It("fails if no endpoint url is given", func() {

			})

			It("defaults files to keep to 5 when it's not provided", func() {
			})
		})
	})
})

func runPipedCmdsAndReturnLastOutput(firstCmd, lastCmd *exec.Cmd) (string, error) {
	var outputBuffer bytes.Buffer
	var err error

	lastCmd.Stdin, err = firstCmd.StdoutPipe()
	Expect(err).ToNot(HaveOccurred())
	lastCmd.Stdout = &outputBuffer

	firstCmd.Start()
	lastCmd.Start()
	err = lastCmd.Wait()

	return outputBuffer.String(), err
}

func s3Bucket(accessKey, secretKey, bucketName string) *s3.Bucket {
	s3Client = testhelpers.BuildGoamzS3(accessKey, secretKey, s3Server.URL())
	return s3Client.Bucket(bucketName)
}
