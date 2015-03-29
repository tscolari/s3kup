package s3_test

import (
	"github.com/tscolari/s3up/s3"

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

			version := s3.NewVersion(key)

			Expect(version.Path).To(Equal("b/a1e53e1d-9b01-cbb99505ac78/0"))
			Expect(version.BackupName).To(Equal("b/a1e53e1d-9b01-cbb99505ac78"))
			Expect(version.Version).To(Equal("0"))
			Expect(version.LastModified).To(Equal("2015-03-29T11:54:42.819+01:00"))
			Expect(version.Size).To(Equal(int64(4)))
		})
	})
})
