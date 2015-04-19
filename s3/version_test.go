package s3_test

import (
	"time"

	"github.com/tscolari/s3kup/s3"

	goamzs3 "github.com/mitchellh/goamz/s3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version", func() {
	Describe("NewVersion", func() {
		It("converts a goamz Key to a version", func() {
			key := goamzs3.Key{
				Key:          "b/a1e53e1d-9b01-cbb99505ac78/0",
				LastModified: "2015-03-29T11:54:42.819+01:00",
				Size:         4,
				ETag:         "\"098f6bcd4621d373cade4e832627b4f6\"",
				StorageClass: "",
			}

			version, err := s3.NewVersion(key)
			Expect(err).ToNot(HaveOccurred())

			Expect(version.Path).To(Equal("b/a1e53e1d-9b01-cbb99505ac78/0"))
			Expect(version.BackupName).To(Equal("b/a1e53e1d-9b01-cbb99505ac78"))
			Expect(version.Version).To(Equal(int64(0)))
			parsedTime, err := time.Parse(time.RFC3339, "2015-03-29T11:54:42.819+01:00")
			Expect(err).ToNot(HaveOccurred())
			Expect(version.LastModified).To(Equal(parsedTime))
			Expect(version.Size).To(Equal(uint64(4)))
		})

		It("returns an error if version number is in a wrong format", func() {
			key := goamzs3.Key{
				Key:          "b/my-bkp.gz/i-am-wrong",
				LastModified: "2015-03-29T11:54:42.819+01:00",
				Size:         4,
				ETag:         "\"098f6bcd4621d373cade4e832627b4f6\"",
				StorageClass: "",
			}

			_, err := s3.NewVersion(key)
			Expect(err).To(MatchError("Remote version 'b/my-bkp.gz/i-am-wrong' can't be parsed"))
		})

		It("returns an error if date is not recognizable", func() {
			key := goamzs3.Key{
				Key:          "b/my-bkp.gz/0",
				LastModified: "03-03-2015T11:54:42.819+01:00",
				Size:         4,
				ETag:         "\"098f6bcd4621d373cade4e832627b4f6\"",
				StorageClass: "",
			}

			_, err := s3.NewVersion(key)
			Expect(err).To(MatchError("Failed to parse the version timestamp. '03-03-2015T11:54:42.819+01:00' was not recognized"))
		})
	})
})
