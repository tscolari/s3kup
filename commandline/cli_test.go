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
			var bucket *s3.Bucket

			BeforeEach(func() {
				bucket = s3Bucket(accessKey, secretKey, bucketName)
				bucket.PutBucket("")
			})

			It("exits with failure if no input is given through a pipe", func() {
				output, err := backupCmd.Output()
				Expect(output).To(MatchRegexp("not using pipeline"))
				Expect(err).To(HaveOccurred())
			})

			It("stores the backup with the correct name", func() {
				_, err := runPipedCmdsAndReturnLastOutput(inputCmd, backupCmd)
				Expect(err).ToNot(HaveOccurred())

				resp, err := bucket.List(backupName, "", "", 100)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(resp.Contents)).To(Equal(1))
			})

			It("keeps only the number of versions specified", func() {
				firstRunInputCmd := exec.Command("echo", "'store my data'")
				firstRunCmd := exec.Command(cli, "-i", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName, "--no-ssl", "-k", "2")
				_, err := runPipedCmdsAndReturnLastOutput(firstRunInputCmd, firstRunCmd)
				Expect(err).ToNot(HaveOccurred())

				secondRunInputCmd := exec.Command("echo", "'store my data'")
				secondRunCmd := exec.Command(cli, "-i", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName, "--no-ssl", "-k", "2")
				_, err = runPipedCmdsAndReturnLastOutput(secondRunInputCmd, secondRunCmd)
				Expect(err).ToNot(HaveOccurred())

				thirdRunInputCmd := exec.Command("echo", "'store my data'")
				thirdRunCmd := exec.Command(cli, "-i", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName, "--no-ssl", "-k", "2")
				_, err = runPipedCmdsAndReturnLastOutput(thirdRunInputCmd, thirdRunCmd)
				Expect(err).ToNot(HaveOccurred())

				resp, err := bucket.List(backupName, "", "", 100)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(resp.Contents)).To(Equal(2))
			})
		})

		Context("when there is invalid or missing args", func() {
			It("fails if no file name is given", func() {
				backupCmd = exec.Command(cli, "-i", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "--no-ssl", "-k", "3")
				output, err := runPipedCmdsAndReturnLastOutput(inputCmd, backupCmd)
				Expect(err).To(HaveOccurred())

				Expect(output).To(MatchRegexp("missing file name argument"))
			})

			It("fails if no bucket name is given", func() {
				backupCmd = exec.Command(cli, "-i", accessKey, "-s", secretKey, "-e", s3EndpointURL, "-n", backupName, "--no-ssl", "-k", "3")
				output, err := runPipedCmdsAndReturnLastOutput(inputCmd, backupCmd)
				Expect(err).To(HaveOccurred())

				Expect(output).To(MatchRegexp("missing bucket name argument"))
			})

			It("fails if no access key is given", func() {
				backupCmd = exec.Command(cli, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName, "--no-ssl", "-k", "3")
				output, err := runPipedCmdsAndReturnLastOutput(inputCmd, backupCmd)
				Expect(err).To(HaveOccurred())

				Expect(output).To(MatchRegexp("missing access key argument"))
			})

			It("fails if no secret key is given", func() {
				backupCmd = exec.Command(cli, "-i", accessKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName, "--no-ssl", "-k", "3")
				output, err := runPipedCmdsAndReturnLastOutput(inputCmd, backupCmd)
				Expect(err).To(HaveOccurred())

				Expect(output).To(MatchRegexp("missing secret key argument"))
			})

			It("fails if versions to keep is equal to zero", func() {
				backupCmd = exec.Command(cli, "-i", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName, "--no-ssl", "-k", "0")
				output, err := runPipedCmdsAndReturnLastOutput(inputCmd, backupCmd)
				Expect(err).To(HaveOccurred())

				Expect(output).To(MatchRegexp("invalid versions to keep. Must be 1 or greater"))
			})

			It("fails if versions to keep is less than zero", func() {
				backupCmd = exec.Command(cli, "-i", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName, "--no-ssl", "-k", "-3")
				output, err := runPipedCmdsAndReturnLastOutput(inputCmd, backupCmd)
				Expect(err).To(HaveOccurred())

				Expect(output).To(MatchRegexp("invalid versions to keep. Must be 1 or greater"))
			})
		})
	})
})

func runPipedCmdsAndReturnLastOutput(firstCmd, lastCmd *exec.Cmd) (string, error) {
	var outputBuffer bytes.Buffer
	var err error

	lastCmd.Stdin, _ = firstCmd.StdoutPipe()
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
