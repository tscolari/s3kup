package list_test

import (
	"errors"
	"time"

	"github.com/tscolari/s3kup/list"
	"github.com/tscolari/s3kup/s3"
	"github.com/tscolari/s3kup/s3/fakeclient"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lister", func() {
	var lister list.Lister
	var s3Client *fakeclient.Client

	BeforeEach(func() {
		s3Client = &fakeclient.Client{}
		lister = list.New(s3Client)
	})

	It("sends the correct request to the s3 client", func() {
		var listCallPath string
		s3Client.ListCall = func(path string) (s3.Versions, error) {
			listCallPath = path
			return nil, nil
		}

		lister.List("my-backup")
		Expect(listCallPath).To(Equal("my-backup"))
	})

	It("forwards the error if s3 client fails", func() {
		s3Client.ListCall = func(path string) (s3.Versions, error) {
			return nil, errors.New("failed here")
		}

		_, err := lister.List("my-backup")
		Expect(err).To(MatchError("failed here"))
	})

	Context("formating", func() {
		var versionOne int64
		var versionTwo int64

		BeforeEach(func() {
			versionOne = time.Now().UnixNano()
			versionTwo = time.Now().UnixNano()
			s3Client.ListCall = func(path string) (s3.Versions, error) {
				return s3.Versions{
					s3.Version{BackupName: path, Version: versionOne},
					s3.Version{BackupName: path, Version: versionTwo},
				}, nil
			}
		})

		It("returns all remote versions", func() {
			versions, err := lister.List("my-backup")
			Expect(err).ToNot(HaveOccurred())
			Expect(len(versions)).To(Equal(2))
		})

		It("returns the remote versions numbers", func() {
			versions, err := lister.List("my-backup")
			Expect(err).ToNot(HaveOccurred())
			Expect(versions).To(Equal(s3.Versions{
				s3.Version{BackupName: "my-backup", Version: versionOne},
				s3.Version{BackupName: "my-backup", Version: versionTwo},
			}))
		})

	})

})
