package integration_test

import (
	"fmt"

	"code.google.com/p/go-uuid/uuid"

	goamzs3 "github.com/mitchellh/goamz/s3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tscolari/s3up/backup"
	"github.com/tscolari/s3up/s3"
	"github.com/tscolari/s3up/s3/testhelpers"
)

var _ = Describe("Backuper", func() {

	const (
		accessKey      string = "my_id"
		secretKey      string = "my_secret"
		regionName     string = "my_region"
		bucketName     string = "my_bucket"
		versionsToKeep int    = 3
	)

	var filePath string
	var backuper backup.Backuper
	var s3Client *goamzs3.S3
	var s3Bucket *goamzs3.Bucket

	BeforeEach(func() {
		client := s3.New(accessKey, secretKey, bucketName, s3EndpointURL)
		backuper = backup.New(client, versionsToKeep)

		s3Client = testhelpers.BuildGoamzS3(accessKey, secretKey, s3EndpointURL)
		s3Bucket = s3Client.Bucket(bucketName)
		s3Bucket.PutBucket("")

		filePath = uuid.New()
	})

	Context("storing the file", func() {
		var data []byte

		BeforeEach(func() {
			data = []byte("file content")
		})

		It("creates a versioned file on s3", func() {
			err := backuper.Backup(filePath, data)
			Expect(err).ToNot(HaveOccurred())

			resp, err := s3Bucket.List(filePath, "", "", 100)
			Expect(err).ToNot(HaveOccurred())

			fileNames := resp.Contents
			Expect(len(fileNames)).To(Equal(1))
			Expect(fileNames[0].Key).To(MatchRegexp(fmt.Sprintf("^%s/\\d{19}$", filePath)))
		})

		It("uploads the correct content to s3", func() {
			err := backuper.Backup(filePath, data)
			Expect(err).ToNot(HaveOccurred())

			resp, err := s3Bucket.List(filePath, "", "", 100)
			Expect(err).ToNot(HaveOccurred())

			fileName := resp.Contents[0].Key
			content, err := s3Bucket.Get(fileName)
			Expect(err).ToNot(HaveOccurred())
			Expect(content).To(Equal(data))
		})
	})

	Context("keeping track of versions", func() {
		BeforeEach(func() {
			for i := 0; i < 3; i++ {
				err := backuper.Backup(filePath, []byte(fmt.Sprintf("data %d", i)))
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("creates different versions for each upload", func() {
			resp, err := s3Bucket.List(filePath, "", "", 100)
			Expect(err).ToNot(HaveOccurred())

			fileNames := []string{}
			for _, fileName := range resp.Contents {
				fileNames = append(fileNames, fileName.Key)
			}
			Expect(len(fileNames)).To(Equal(3))
		})

		It("keeps the correct content on each version", func() {
			resp, err := s3Bucket.List(filePath, "", "", 100)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(resp.Contents)).To(Equal(3))

			for i, fileName := range resp.Contents {
				content, err := s3Bucket.Get(fileName.Key)
				Expect(err).ToNot(HaveOccurred())
				Expect(content).To(Equal([]byte(fmt.Sprintf("data %d", i))))
			}
		})
	})

	Context("cleaning up old versions", func() {
		BeforeEach(func() {
			for i := 0; i < 5; i++ {
				fileName := fmt.Sprintf("%s/0000000000000000%d", filePath, i)
				s3Bucket.Put(fileName, []byte(""), "", "")
			}
		})

		It("cleans up the oldest versions until only `versionsToKeep` versions exist", func() {
			resp, err := s3Bucket.List(filePath, "", "", 100)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(resp.Contents)).To(Equal(5))

			err = backuper.Backup(filePath, []byte("data"))
			Expect(err).ToNot(HaveOccurred())

			resp, err = s3Bucket.List(filePath, "", "", 100)
			Expect(err).ToNot(HaveOccurred())
			fileNames := []string{}
			for _, fileName := range resp.Contents {
				fileNames = append(fileNames, fileName.Key)
			}
			Expect(len(fileNames)).To(Equal(3))

			Expect(fileNames).ToNot(ContainElement(fmt.Sprintf("%s/00000000000000000", filePath)))
			Expect(fileNames).ToNot(ContainElement(fmt.Sprintf("%s/00000000000000001", filePath)))
			Expect(fileNames).ToNot(ContainElement(fmt.Sprintf("%s/00000000000000002", filePath)))

			Expect(fileNames).To(ContainElement(fmt.Sprintf("%s/00000000000000003", filePath)))
		})

	})
})
