package integration_test

import (
	"github.com/google/uuid"
	goamzs3 "github.com/mitchellh/goamz/s3"
	"github.com/tscolari/s3kup/list"
	"github.com/tscolari/s3kup/s3"
	"github.com/tscolari/s3kup/s3/testhelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lister", func() {

	const (
		accessKey      string = "my_id"
		secretKey      string = "my_secret"
		regionName     string = "my_region"
		versionsToKeep int    = 3
	)

	var filePath string
	var lister list.Lister
	var s3Client *goamzs3.S3
	var s3Bucket *goamzs3.Bucket

	BeforeEach(func() {
		bucketName := uuid.New()
		client := s3.New(accessKey, secretKey, bucketName, s3EndpointURL)
		lister = list.New(client)

		s3Client = testhelpers.BuildGoamzS3(accessKey, secretKey, s3EndpointURL)
		s3Bucket = s3Client.Bucket(bucketName)

		filePath = uuid.New()
	})

	Context("when things work", func() {
		BeforeEach(func() {
			s3Bucket.PutBucket("")
		})

		Context("when there are 0 versions stored", func() {
			It("returns an empty Versions object", func() {
				versions, err := lister.List("noversions")
				Expect(err).ToNot(HaveOccurred())

				Expect(len(versions)).To(Equal(0))
			})
		})

		Context("when there are n versions stored", func() {
			BeforeEach(func() {
				s3Bucket.Put("my-db/100001", []byte("data"), "", "")
				s3Bucket.Put("my-db/100003", []byte("data"), "", "")
				s3Bucket.Put("my-db/100002", []byte("data"), "", "")
				s3Bucket.Put("my-db/100004", []byte("data"), "", "")
			})

			It("returns a sorted Versions object with all versions", func() {
				versions, err := lister.List("my-db")
				Expect(err).ToNot(HaveOccurred())

				Expect(len(versions)).To(Equal(4))
				Expect(versions[0].Path).To(Equal("my-db/100001"))
				Expect(versions[1].Path).To(Equal("my-db/100002"))
				Expect(versions[2].Path).To(Equal("my-db/100003"))
				Expect(versions[3].Path).To(Equal("my-db/100004"))
			})
		})
	})

	Context("when there is no bucket", func() {
		It("fails to list the versions", func() {
			_, err := lister.List("my-db")
			Expect(err).To(MatchError("The specified bucket does not exist"))
		})
	})
})
