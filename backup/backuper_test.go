package backup_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/tscolari/s3kup/backup"
	"github.com/tscolari/s3kup/backup/fakes"
	"github.com/tscolari/s3kup/s3"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Backuper", func() {
	var backuper backup.Backuper
	var s3Client *fakes.FakeS3Client

	BeforeEach(func() {
		s3Client = new(fakes.FakeS3Client)
		backuper = backup.New(s3Client, 3)
	})

	Describe("#Backup", func() {

		It("timestamps the version inside the given filename", func() {
			err := backuper.Backup("file", []byte("content"))
			Expect(err).ToNot(HaveOccurred())
			path, _ := s3Client.StoreArgsForCall(0)
			Expect(path).To(MatchRegexp(fmt.Sprintf("^%s/\\d{19}$", "file")))
		})

		Context("when something fails", func() {

			Context("when storing the file fails", func() {
				It("returns back the error", func() {
					s3Client.StoreReturns(errors.New("failed to store"))
					err := backuper.Backup("file", []byte("content"))
					Expect(err).To(MatchError("failed to store"))
				})
			})

			Context("when listing the versions fails", func() {

				BeforeEach(func() {
					s3Client.StoreReturns(errors.New("Failed to list"))
				})

				It("returns back the error", func() {
					err := backuper.Backup("file", []byte("content"))
					Expect(err).To(MatchError("Failed to list"))
				})

				It("still stores the file", func() {
					backuper.Backup("file", []byte("content"))
					Expect(s3Client.StoreCallCount()).To(Equal(1))
				})
			})

			Context("when deleting old versions fails", func() {
				BeforeEach(func() {
					s3Client.DeleteReturns(errors.New("Failed to delete"))
					versions := s3.Versions{
						s3.Version{BackupName: "myfile", Version: 20010101},
						s3.Version{BackupName: "myfile", Version: 20000101},
						s3.Version{BackupName: "myfile", Version: 20020101},
						s3.Version{BackupName: "myfile", Version: 20150101},
					}

					s3Client.ListReturns(versions, nil)
				})

				It("returns back the error", func() {
					err := backuper.Backup("file", []byte("content"))
					Expect(err).To(MatchError("Failed to delete"))
				})

				It("still store the file", func() {
					backuper.Backup("file", []byte("content"))
					Expect(s3Client.StoreCallCount()).To(Equal(1))
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
					s3Client.ListReturns(s3.Versions{}, nil)
					s3Client.DeleteReturns(nil)

					err := backuper.Backup("file", []byte("content"))
					Expect(err).ToNot(HaveOccurred())
					Expect(s3Client.DeleteCallCount()).To(Equal(0))
				})
			})

			Context("when there is as many versions as `versionsToKeep`", func() {
				It("deletes the oldest version", func() {
					baseTime := time.Now()
					versions := s3.Versions{
						s3.Version{BackupName: "myfile", Version: 20010101, Path: "myfile/20010101", LastModified: baseTime.Add(2 * time.Minute)},
						s3.Version{BackupName: "myfile", Version: 20000101, Path: "myfile/20000101", LastModified: baseTime.Add(1 * time.Minute)},
						s3.Version{BackupName: "myfile", Version: 20020101, Path: "myfile/20020101", LastModified: baseTime.Add(3 * time.Minute)},
						s3.Version{BackupName: "myfile", Version: 20150101, Path: "myfile/20150101", LastModified: baseTime.Add(4 * time.Minute)},
					}

					s3Client.ListReturns(versions, nil)

					err := backuper.Backup("file", []byte("content"))
					Expect(err).ToNot(HaveOccurred())
					Expect(s3Client.DeleteCallCount()).To(Equal(1))
					deletedPath := s3Client.DeleteArgsForCall(0)
					Expect(deletedPath).To(Equal("myfile/20000101"))
				})
			})

			Context("when there is more versions than `versionsToKeep`", func() {
				It("deletes as many old versions as necessary to keep it the same as `versionsToKeep`", func() {
					baseTime := time.Now()

					versions := s3.Versions{
						s3.Version{BackupName: "myfile", Version: 20010101, Path: "myfile/20010101", LastModified: baseTime.Add(4 * time.Minute)},
						s3.Version{BackupName: "myfile", Version: 20030101, Path: "myfile/20030101", LastModified: baseTime.Add(6 * time.Minute)},
						s3.Version{BackupName: "myfile", Version: 20000101, Path: "myfile/20000101", LastModified: baseTime.Add(3 * time.Minute)},
						s3.Version{BackupName: "myfile", Version: 20020101, Path: "myfile/20020101", LastModified: baseTime.Add(5 * time.Minute)},
						s3.Version{BackupName: "myfile", Version: 19990101, Path: "myfile/19990101", LastModified: baseTime.Add(2 * time.Minute)},
						s3.Version{BackupName: "myfile", Version: 20150101, Path: "myfile/20150101", LastModified: baseTime.Add(7 * time.Minute)},
						s3.Version{BackupName: "myfile", Version: 19950101, Path: "myfile/19950101", LastModified: baseTime.Add(1 * time.Minute)},
					}

					s3Client.ListReturns(versions, nil)

					err := backuper.Backup("file", []byte("content"))
					Expect(err).ToNot(HaveOccurred())
					Expect(s3Client.DeleteCallCount()).To(Equal(4))
					Expect(s3Client.DeleteArgsForCall(0)).To(Equal("myfile/19950101"))
					Expect(s3Client.DeleteArgsForCall(1)).To(Equal("myfile/19990101"))
					Expect(s3Client.DeleteArgsForCall(2)).To(Equal("myfile/20000101"))
					Expect(s3Client.DeleteArgsForCall(3)).To(Equal("myfile/20010101"))
				})
			})
		})
	})
})
