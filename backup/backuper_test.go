package backup_test

import (
	"errors"
	"fmt"

	"github.com/tscolari/s3up/backup"
	"github.com/tscolari/s3up/s3"
	"github.com/tscolari/s3up/s3/fakeclient"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Backuper", func() {
	var backuper backup.Backuper
	var s3Client *fakeclient.Client

	BeforeEach(func() {
		s3Client = &fakeclient.Client{}
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
			Expect(receivedFileName).To(MatchRegexp(fmt.Sprintf("^%s/\\d{19}$", "file")))
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
					s3Client.ListCall = func(path string) (s3.Versions, error) {
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

					s3Client.ListCall = func(path string) (s3.Versions, error) {
						versions := s3.Versions{
							s3.Version{BackupName: "myfile", Version: 20010101},
							s3.Version{BackupName: "myfile", Version: 20000101},
							s3.Version{BackupName: "myfile", Version: 20020101},
							s3.Version{BackupName: "myfile", Version: 20150101},
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
					s3Client.ListCall = func(path string) (s3.Versions, error) {
						return s3.Versions{}, nil
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

					s3Client.ListCall = func(path string) (s3.Versions, error) {
						versions := s3.Versions{
							s3.Version{BackupName: "myfile", Version: 20010101, Path: "myfile/20010101"},
							s3.Version{BackupName: "myfile", Version: 20000101, Path: "myfile/20000101"},
							s3.Version{BackupName: "myfile", Version: 20020101, Path: "myfile/20020101"},
							s3.Version{BackupName: "myfile", Version: 20150101, Path: "myfile/20150101"},
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
					Expect(deletedVersion).To(Equal("myfile/20000101"))
				})
			})

			Context("when there is more versions than `versionsToKeep`", func() {
				It("deletes as many old versions as necessary to keep it the same as `versionsToKeep`", func() {
					deletedVersions := []string{}

					s3Client.ListCall = func(path string) (s3.Versions, error) {
						versions := s3.Versions{
							s3.Version{BackupName: "myfile", Version: 20010101, Path: "myfile/20010101"},
							s3.Version{BackupName: "myfile", Version: 20030101, Path: "myfile/20030101"},
							s3.Version{BackupName: "myfile", Version: 20000101, Path: "myfile/20000101"},
							s3.Version{BackupName: "myfile", Version: 20020101, Path: "myfile/20020101"},
							s3.Version{BackupName: "myfile", Version: 19990101, Path: "myfile/19990101"},
							s3.Version{BackupName: "myfile", Version: 20150101, Path: "myfile/20150101"},
							s3.Version{BackupName: "myfile", Version: 19950101, Path: "myfile/19950101"},
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
						"myfile/19950101",
						"myfile/19990101",
						"myfile/20000101",
						"myfile/20010101",
					}))
				})
			})
		})
	})
})
