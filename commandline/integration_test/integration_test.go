package integration_test

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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
		backupCmd = exec.Command(cli, "-i", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName)

		bucketName = fmt.Sprintf("bucket%d", rand.Int())
		bucket := s3Bucket(accessKey, secretKey, bucketName)
		bucket.PutBucket("")
	})

	Context("when there is invalid or missing args", func() {
		It("fails if no bucket name is given", func() {
			backupCmd = exec.Command(cli, "-i", accessKey, "-s", secretKey, "-e", s3EndpointURL, "-n", backupName)
			output, err := runPipedCmdsAndReturnLastOutput(inputCmd, backupCmd)
			Expect(err).To(HaveOccurred())

			Expect(output).To(MatchRegexp("missing bucket name argument"))
		})

		It("fails if no access key is given", func() {
			backupCmd = exec.Command(cli, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName)
			output, err := runPipedCmdsAndReturnLastOutput(inputCmd, backupCmd)
			Expect(err).To(HaveOccurred())

			Expect(output).To(MatchRegexp("missing access key argument"))
		})

		It("fails if no secret key is given", func() {
			backupCmd = exec.Command(cli, "-i", accessKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName)
			output, err := runPipedCmdsAndReturnLastOutput(inputCmd, backupCmd)
			Expect(err).To(HaveOccurred())

			Expect(output).To(MatchRegexp("missing secret key argument"))
		})
	})

	Describe("uploading and fetching the same file", func() {
		var pushCmd *exec.Cmd
		var inputCmd *exec.Cmd
		var pullCmd *exec.Cmd
		var unzipCmd *exec.Cmd
		var inputData []byte

		BeforeEach(func() {
			pushCmd = exec.Command(cli, "push", "-i", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName)
			pullCmd = exec.Command(cli, "pull", "-i", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName)
		})

		Context("Text files", func() {
			BeforeEach(func() {
				var err error
				inputData, err = ioutil.ReadFile("./assets/backupfile")
				Expect(err).ToNot(HaveOccurred())
				inputCmd = exec.Command("cat", "./assets/backupfile")
			})

			It("should pull and output exactly the same content that was pushed", func() {
				_, err := runPipedCmdsAndReturnLastOutput(inputCmd, pushCmd)
				Expect(err).ToNot(HaveOccurred())

				output, err := pullCmd.CombinedOutput()
				Expect(err).ToNot(HaveOccurred())
				Expect(output).To(Equal(inputData))
			})
		})

		Context("Binary files", func() {
			BeforeEach(func() {
				var err error
				inputData, err = ioutil.ReadFile("./assets/gopher.png")
				Expect(err).ToNot(HaveOccurred())
				inputCmd = exec.Command("cat", "./assets/gopher.png")
			})

			It("should pull and output exactly the same content that was pushed", func() {
				_, err := runPipedCmdsAndReturnLastOutput(inputCmd, pushCmd)
				Expect(err).ToNot(HaveOccurred())

				output, err := pullCmd.CombinedOutput()
				Expect(err).ToNot(HaveOccurred())
				Expect(md5.Sum(output)).To(Equal(md5.Sum(inputData)))
			})

			Context("Using compression", func() {

				BeforeEach(func() {
					inputCmd = exec.Command("gzip", "./assets/gopher.png", "-c")
					unzipCmd = exec.Command("gunzip")
				})

				It("should pull and output exactly the same content that was pushed", func() {
					_, err := runPipedCmdsAndReturnLastOutput(inputCmd, pushCmd)
					Expect(err).ToNot(HaveOccurred())

					output, err := runPipedCmdsAndReturnLastOutput(pullCmd, unzipCmd)
					Expect(err).ToNot(HaveOccurred())
					Expect(md5.Sum([]byte(output))).To(Equal(md5.Sum(inputData)))
				})
			})
		})
	})
})
