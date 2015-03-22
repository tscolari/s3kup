package backup_test

import (
	"errors"
	"fmt"

	"github.com/s3up/backup"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakeS3Client struct {
	StoreCall  func(path string, content []byte) error
	ListCall   func(path string) ([]string, error)
	DeleteCall func(path string) error
}

func (c *fakeS3Client) Store(path string, content []byte) error {
	if c.StoreCall != nil {
		return c.StoreCall(path, content)
	}
	return nil
}

func (c *fakeS3Client) List(path string) (files []string, err error) {
	if c.ListCall != nil {
		return c.ListCall(path)
	}
	return nil, nil
}

func (c *fakeS3Client) Delete(path string) error {
	if c.DeleteCall != nil {
		return c.DeleteCall(path)
	}
	return nil
}

var _ = Describe("Backuper", func() {
	var backuper backup.Backuper
	var s3Client *fakeS3Client

	BeforeEach(func() {
		s3Client = &fakeS3Client{}
		backuper = backup.New(s3Client, 3)
	})

	Describe("#Backup", func() {
		var receivedFileName string
		var receivedContent []byte

		BeforeEach(func() {
			s3Client.StoreCall = func(path string, content []byte) error {
				receivedFileName = path
				receivedContent = content
				return nil
			}
		})

		It("timestamps the version inside the given filename", func() {
			err := backuper.Backup("file", []byte("content"))
			Expect(err).ToNot(HaveOccurred())
			Expect(receivedFileName).To(MatchRegexp(fmt.Sprintf("%s/\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}", "file")))
		})

		Context("when something fails", func() {
			var storeCallCount int

			BeforeEach(func() {
				storeCallCount = 0
				s3Client.StoreCall = func(path string, content []byte) error {
					storeCallCount++
					return nil
				}
			})

			Context("when storing the file fails", func() {
				It("returns back the error", func() {
					s3Client.StoreCall = func(path string, content []byte) error {
						return errors.New("failed to store")
					}
					err := backuper.Backup("file", []byte("content"))
					Expect(err).To(MatchError("failed to store"))
				})
			})

			Context("when listing the versions fails", func() {

				BeforeEach(func() {
					s3Client.ListCall = func(path string) ([]string, error) {
						return nil, errors.New("Failed to list")
					}
				})

				It("returns back the error", func() {
					err := backuper.Backup("file", []byte("content"))
					Expect(err).To(MatchError("Failed to list"))
				})

				It("still stores the file", func() {
					backuper.Backup("file", []byte("content"))
					Expect(storeCallCount).To(Equal(1))
				})
			})

			Context("when deleting old versions fails", func() {
				BeforeEach(func() {
					s3Client.DeleteCall = func(path string) error {
						return errors.New("Failed to delete")
					}

					s3Client.ListCall = func(path string) ([]string, error) {
						versions := []string{
							"myfile/2001-01-01",
							"myfile/2000-01-01",
							"myfile/2002-01-01",
							"myfile/2015-01-01-just-pushed",
						}
						return versions, nil
					}
				})

				It("returns back the error", func() {
					err := backuper.Backup("file", []byte("content"))
					Expect(err).To(MatchError("Failed to delete"))
				})

				It("still store the file", func() {
					backuper.Backup("file", []byte("content"))
					Expect(storeCallCount).To(Equal(1))
				})
			})
		})

		Context("versions to keep", func() {
			var storeCallsCount int
			var deleteCallsCount int
			var listCallsCount int

			BeforeEach(func() {
				storeCallsCount = 0
				deleteCallsCount = 0
				listCallsCount = 0
			})

			Context("when there is less versions than `versionsToKeep`", func() {
				It("does not delete any previous version", func() {
					s3Client.ListCall = func(path string) ([]string, error) {
						return []string{}, nil
					}

					s3Client.DeleteCall = func(path string) error {
						deleteCallsCount++
						return nil
					}

					err := backuper.Backup("file", []byte("content"))
					Expect(err).ToNot(HaveOccurred())
					Expect(deleteCallsCount).To(Equal(0))
				})
			})

			Context("when there is as many versions as `versionsToKeep`", func() {
				It("deletes the oldest version", func() {
					var deletedVersion string

					s3Client.ListCall = func(path string) ([]string, error) {
						versions := []string{
							"myfile/2001-01-01",
							"myfile/2000-01-01",
							"myfile/2002-01-01",
							"myfile/2015-01-01-just-pushed",
						}
						return versions, nil
					}

					s3Client.DeleteCall = func(path string) error {
						deletedVersion = path
						deleteCallsCount++
						return nil
					}

					err := backuper.Backup("file", []byte("content"))
					Expect(err).ToNot(HaveOccurred())
					Expect(deleteCallsCount).To(Equal(1))
					Expect(deletedVersion).To(Equal("myfile/2000-01-01"))
				})
			})

			Context("when there is more versions than `versionsToKeep`", func() {
				It("deletes as many old versions as necessary to keep it the same as `versionsToKeep`", func() {
					deletedVersions := []string{}

					s3Client.ListCall = func(path string) ([]string, error) {
						versions := []string{
							"myfile/2001-01-01",
							"myfile/2003-01-01",
							"myfile/2000-01-01",
							"myfile/2002-01-01",
							"myfile/1999-01-01",
							"myfile/2015-01-01-just-pushed",
							"myfile/1995-01-01",
						}
						return versions, nil
					}

					s3Client.DeleteCall = func(path string) error {
						deletedVersions = append(deletedVersions, path)
						deleteCallsCount++
						return nil
					}

					err := backuper.Backup("file", []byte("content"))
					Expect(err).ToNot(HaveOccurred())
					Expect(deleteCallsCount).To(Equal(4))
					Expect(deletedVersions).To(Equal([]string{
						"myfile/1995-01-01",
						"myfile/1999-01-01",
						"myfile/2000-01-01",
						"myfile/2001-01-01",
					}))
				})
			})
		})
	})
})
