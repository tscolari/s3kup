package integration_test

import (
	"code.google.com/p/go-uuid/uuid"
	goamzs3 "github.com/mitchellh/goamz/s3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tscolari/s3up/fetch"
	"github.com/tscolari/s3up/s3"
	"github.com/tscolari/s3up/s3/testhelpers"
)

var _ = Describe("Fetcher", func() {

	const (
		accessKey  string = "my_id"
		secretKey  string = "my_secret"
		regionName string = "my_region"
	)

	var filePath string
	var fetcher fetch.Fetcher
	var s3Client *goamzs3.S3
	var s3Bucket *goamzs3.Bucket
	var backupName string

	BeforeEach(func() {
		backupName = "my/backup"
		bucketName := uuid.New()
		client := s3.New(accessKey, secretKey, bucketName, s3EndpointURL)
		fetcher = fetch.New(client)

		s3Client = testhelpers.BuildGoamzS3(accessKey, secretKey, s3EndpointURL)
		s3Bucket = s3Client.Bucket(bucketName)
		s3Bucket.PutBucket("")

		s3Bucket.Put(backupName+"/1001", []byte("first backup"), "", "")
		s3Bucket.Put(backupName+"/1003", []byte("third backup"), "", "")
		s3Bucket.Put(backupName+"/1002", []byte("second backup"), "", "")

		filePath = uuid.New()
	})

	Describe("#FetchLatest", func() {

		Context("when the bucket is empty", func() {
			It("returns an error", func() {
				_, err := fetcher.FetchLatest("do/not/exist")
				Expect(err).To(MatchError("There's no backup named 'do/not/exist' on this bucket"))
			})
		})

		Context("when there are versions stored", func() {
			It("returns the content of the latest version", func() {
				content, err := fetcher.FetchLatest(backupName)
				Expect(err).ToNot(HaveOccurred())
				Expect(content).To(Equal([]byte("third backup")))
			})
		})
	})

	Describe("#FetchVersion", func() {
		Context("when the version doesn't exist", func() {
			It("returns an error", func() {
				_, err := fetcher.FetchVersion(backupName, 999)
				Expect(err).To(MatchError("Could not find version '999'"))
			})
		})

		It("returns the content of the given version", func() {
			content, err := fetcher.FetchVersion(backupName, 1001)
			Expect(err).ToNot(HaveOccurred())

			Expect(content).To(Equal([]byte("first backup")))
		})

	})
})
