package s3_test

import (
	"fmt"
	"math/rand"

	"code.google.com/p/go-uuid/uuid"

	"github.com/tscolari/s3up/s3"
	"github.com/tscolari/s3up/s3/testhelpers"

	goamzs3 "github.com/mitchellh/goamz/s3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	const (
		accessKey  string = "my_id"
		secretKey  string = "my_secret"
		regionName string = "my_region"
		bucketName string = "my_bucket"
	)

	var client *s3.Client
	var bucket *goamzs3.Bucket
	var filePath string

	BeforeEach(func() {
		filePath = uuid.New()
		s3Client := testhelpers.BuildGoamzS3(accessKey, secretKey, s3EndpointURL)

		bucket = s3Client.Bucket(bucketName)
		bucket.PutBucket(goamzs3.Private)

		client = s3.New(accessKey, secretKey, bucketName, s3EndpointURL)
	})

	AfterEach(func() {
		bucket.DelBucket()
	})

	Describe("#Store", func() {
		It("stores the file content with the file name on s3", func() {
			fileContent := []byte("my file contents")

			err := client.Store(filePath, fileContent)
			Expect(err).ToNot(HaveOccurred())

			remoteContent, err := bucket.Get(filePath)
			Expect(err).ToNot(HaveOccurred())
			Expect(remoteContent).To(Equal(fileContent))
		})
	})

	Describe("#List", func() {
		BeforeEach(func() {
			for i := 4; i >= 0; i-- {
				path := fmt.Sprintf("%s/%d", filePath, i)
				err := bucket.Put(path, []byte("test"), "", "")
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("lists all files, sorted, in the s3 bucket path", func() {
			files, err := client.List(filePath)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(files)).To(Equal(5))

			for i := 0; i < 5; i++ {
				Expect(files[i].Version).To(Equal(int64(i)))
				Expect(files[i].BackupName).To(Equal(filePath))
			}
		})
	})

	Describe("#Get", func() {
		var path string
		var remoteContent []byte

		BeforeEach(func() {
			path = fmt.Sprintf("%s/%d", rand.Int())
			remoteContent = []byte(fmt.Sprintf("content %d", rand.Int()))

			err := bucket.Put(path, remoteContent, "", "")
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns the correct content", func() {
			content, err := client.Get(path)
			Expect(err).ToNot(HaveOccurred())
			Expect(content).To(Equal(remoteContent))
		})
	})

	Describe("#Delete", func() {
		It("removes the s3 file path", func() {
			err := bucket.Put(filePath, []byte("test"), "", "")
			Expect(err).ToNot(HaveOccurred())

			err = client.Delete(filePath)
			Expect(err).ToNot(HaveOccurred())

			_, err = bucket.Get(filePath)
			Expect(err).To(MatchError("The specified key does not exist."))
		})
	})
})
