package commandline_test

import (
	"fmt"
	"math/rand"
	"os/exec"

	"github.com/mitchellh/goamz/s3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cli > list", func() {

	const (
		accessKey  string = "my_id"
		secretKey  string = "my_secret"
		regionName string = "my_region"
		backupName string = "my/backup"
	)

	var listCmd *exec.Cmd
	var bucket *s3.Bucket

	BeforeEach(func() {
		bucketName := fmt.Sprintf("bucket%d", rand.Int())
		listCmd = exec.Command(cli, "list", "-i", accessKey, "-s", secretKey, "-b", bucketName, "-e", s3EndpointURL, "-n", backupName)
		bucket = s3Bucket(accessKey, secretKey, bucketName)
		bucket.PutBucket("")
	})

	Context("when there are remote versions", func() {
		BeforeEach(func() {
			bucket.Put("my/backup/10000001", []byte("content"), "", "")
			bucket.Put("my/backup/10000002", []byte("content"), "", "")
			bucket.Put("my/backup/10000003", []byte("content"), "", "")
			bucket.Put("my/backup/10000004", []byte("content"), "", "")
		})

		It("lists all remote files", func() {
			output, err := listCmd.CombinedOutput()
			Expect(err).ToNot(HaveOccurred())

			Expect(string(output)).To(MatchRegexp("\\* 10000001"))
			Expect(string(output)).To(MatchRegexp("\\* 10000002"))
			Expect(string(output)).To(MatchRegexp("\\* 10000003"))
			Expect(string(output)).To(MatchRegexp("\\* 10000004"))
		})
	})

	Context("when there aren't remote versions", func() {
		It("prints a no versions stored message", func() {
			output, err := listCmd.CombinedOutput()
			Expect(err).ToNot(HaveOccurred())

			Expect(string(output)).To(MatchRegexp("No versions found"))
		})
	})
})
