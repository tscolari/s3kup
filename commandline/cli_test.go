package commandline_test

import (
	"fmt"
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
})
