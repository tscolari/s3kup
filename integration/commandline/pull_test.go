package integration_test

import (
	"fmt"
	"math/rand"
	"os/exec"

	"github.com/mitchellh/goamz/s3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cli > pull", func() {

	const (
		accessKey  string = "my_id"
		secretKey  string = "my_secret"
		regionName string = "my_region"
		backupName string = "my/backup"
	)

	var pullCmd *exec.Cmd
	var bucket *s3.Bucket
	var bucketName string

	BeforeEach(func() {
		bucketName = fmt.Sprintf("bucket%d", rand.Int())
		pullCmd = exec.Command(cli, "pull", "-a", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName)
		bucket = s3Bucket(accessKey, secretKey, bucketName)
		bucket.PutBucket("")

	})

	Context("when no argument is given", func() {
		Context("when there are remote versions", func() {
			BeforeEach(func() {
				bucket.Put("my/backup/10000002", []byte("content 2"), "", "")
				bucket.Put("my/backup/10000001", []byte("content 1"), "", "")
				bucket.Put("my/backup/10000004", []byte("content 4"), "", "")
				bucket.Put("my/backup/10000003", []byte("content 3"), "", "")
			})

			It("fetches and prints out the latest backup", func() {
				output, err := pullCmd.CombinedOutput()
				Expect(err).ToNot(HaveOccurred())
				Expect(output).To(Equal([]byte("content 4")))
			})
		})

		Context("when there is no backups", func() {
			It("prints an error", func() {
				output, err := pullCmd.CombinedOutput()
				Expect(err).To(HaveOccurred())

				Expect(string(output)).To(MatchRegexp("There's no backup named 'my/backup' on this bucket"))
			})
		})
	})

	Context("when a version is specified", func() {
		BeforeEach(func() {
			bucket.Put("my/backup/10000002", []byte("content 2"), "", "")
			bucket.Put("my/backup/10000001", []byte("content 1"), "", "")
			bucket.Put("my/backup/10000004", []byte("content 4"), "", "")
			bucket.Put("my/backup/10000003", []byte("content 3"), "", "")
		})

		It("fetches that version's content", func() {
			pullCmd = exec.Command(cli, "pull", "10000003", "-a", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName)
			output, err := pullCmd.CombinedOutput()
			Expect(err).ToNot(HaveOccurred())

			Expect(output).To(Equal([]byte("content 3")))
		})

		Context("when the version doesn't exist", func() {
			It("prints an error", func() {
				pullCmd = exec.Command(cli, "pull", "19999993", "-a", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName)
				output, err := pullCmd.CombinedOutput()
				Expect(err).To(HaveOccurred())

				Expect(string(output)).To(MatchRegexp("Could not find version '19999993'"))
			})
		})
	})
})
