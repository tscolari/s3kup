package integration_test

import (
	"fmt"
	"math/rand"
	"os/exec"

	"github.com/mitchellh/goamz/s3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cli > push", func() {

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
		var bucket *s3.Bucket

		BeforeEach(func() {
			backupCmd = exec.Command(cli, "push", "-a", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName, "-k", "3")
			bucket = s3Bucket(accessKey, secretKey, bucketName)
			bucket.PutBucket("")
		})

		Context("when the correct args are given", func() {

			It("exits with failure if no input is given through a pipe", func() {
				output, err := backupCmd.CombinedOutput()
				Expect(output).To(MatchRegexp("not using pipeline"))
				Expect(err).To(HaveOccurred())
			})

			It("stores the backup with the correct name", func() {
				output, err := runPipedCmdsAndReturnLastOutput(inputCmd, backupCmd)
				fmt.Println(output)
				Expect(err).ToNot(HaveOccurred())

				resp, err := bucket.List(backupName, "", "", 100)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(resp.Contents)).To(Equal(1))
			})

			It("keeps only the number of versions specified", func() {
				firstRunInputCmd := exec.Command("echo", "'store my data'")
				firstRunCmd := exec.Command(cli, "push", "-a", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName, "-k", "2")
				_, err := runPipedCmdsAndReturnLastOutput(firstRunInputCmd, firstRunCmd)
				Expect(err).ToNot(HaveOccurred())

				secondRunInputCmd := exec.Command("echo", "'store my data'")
				secondRunCmd := exec.Command(cli, "push", "-a", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName, "-k", "2")
				_, err = runPipedCmdsAndReturnLastOutput(secondRunInputCmd, secondRunCmd)
				Expect(err).ToNot(HaveOccurred())

				thirdRunInputCmd := exec.Command("echo", "'store my data'")
				thirdRunCmd := exec.Command(cli, "push", "-a", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName, "-k", "2")
				_, err = runPipedCmdsAndReturnLastOutput(thirdRunInputCmd, thirdRunCmd)
				Expect(err).ToNot(HaveOccurred())

				resp, err := bucket.List(backupName, "", "", 100)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(resp.Contents)).To(Equal(2))
			})

			Context("on verbose mode", func() {
				It("outputs the steps", func() {
					backupCmd := exec.Command(cli, "push", "-a", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName, "-k", "3", "--verbose")
					output, err := runPipedCmdsAndReturnLastOutput(inputCmd, backupCmd)
					Expect(err).ToNot(HaveOccurred())

					Expect(output).To(MatchRegexp("Started backup of my-backup"))
					Expect(output).To(MatchRegexp(" -- File version: \\d{19}\\]"))
					Expect(output).To(MatchRegexp(" -- Looking for old versions to delete. keeping 3"))
				})
			})
		})

		Context("when there is invalid or missing args", func() {
			It("fails if no file name is given", func() {
				backupCmd = exec.Command(cli, "push", "-a", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-k", "3")
				output, err := runPipedCmdsAndReturnLastOutput(inputCmd, backupCmd)
				Expect(err).To(HaveOccurred())

				Expect(output).To(MatchRegexp("missing file name argument"))
			})

			It("fails if versions to keep is equal to zero", func() {
				backupCmd = exec.Command(cli, "push", "-a", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName, "-k", "0")
				output, err := runPipedCmdsAndReturnLastOutput(inputCmd, backupCmd)
				Expect(err).To(HaveOccurred())

				Expect(output).To(MatchRegexp("invalid versions to keep. Must be 1 or greater"))
			})

			It("fails if versions to keep is less than zero", func() {
				backupCmd = exec.Command(cli, "push", "-a", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName, "-k", "-3")
				output, err := runPipedCmdsAndReturnLastOutput(inputCmd, backupCmd)
				Expect(err).To(HaveOccurred())

				Expect(output).To(MatchRegexp("invalid versions to keep. Must be 1 or greater"))
			})
		})
	})
})
