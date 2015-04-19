package fetch_test

import (
	"errors"

	"github.com/tscolari/s3kup/fetch"
	"github.com/tscolari/s3kup/s3"
	"github.com/tscolari/s3kup/s3/fakeclient"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fetcher", func() {
	var fetcher fetch.Fetcher
	var client *fakeclient.Client

	BeforeEach(func() {
		client = &fakeclient.Client{
			ListCall: func(path string) (s3.Versions, error) {
				return s3.Versions{
					s3.Version{Path: "my-backup/0", Version: 0},
					s3.Version{Path: "my-backup/1", Version: 1},
					s3.Version{Path: "my-backup/2", Version: 2},
				}, nil
			},
		}
		fetcher = fetch.New(client)
	})

	Describe("#FetchLatest", func() {
		Context("when the s3 client returns an error", func() {
			It("forwards the error", func() {
				client.GetCall = func(path string) ([]byte, error) {
					return nil, errors.New("some error")
				}

				_, err := fetcher.FetchLatest("my-backup")
				Expect(err).To(MatchError("some error"))
			})
		})

		Context("when there is no versions/backup", func() {
			It("returns an error", func() {
				client.ListCall = func(path string) (s3.Versions, error) {
					return s3.Versions{}, nil
				}
				_, err := fetcher.FetchLatest("dontexist")
				Expect(err).To(MatchError("There's no backup named 'dontexist' on this bucket"))
			})
		})

		It("returns the content of the latest version", func() {
			client.GetCall = func(path string) ([]byte, error) {
				if path == "my-backup/2" {
					return []byte("correct version"), nil
				}
				return nil, errors.New("incorrect version")
			}

			content, err := fetcher.FetchLatest("my-backup")
			Expect(err).ToNot(HaveOccurred())
			Expect(content).To(Equal([]byte("correct version")))
		})
	})

	Describe("#FetchVersion", func() {
		Context("when the s3 client returns an error", func() {
			It("forwards the error", func() {
				client.GetCall = func(path string) ([]byte, error) {
					return nil, errors.New("some error")
				}

				_, err := fetcher.FetchVersion("my-backup", 1)
				Expect(err).To(MatchError("some error"))
			})
		})

		It("returns the content of the given version", func() {
			client.GetCall = func(path string) ([]byte, error) {
				if path == "my-backup/1" {
					return []byte("version 1 content"), nil
				}
				return []byte("another version content"), nil
			}

			content, err := fetcher.FetchVersion("my-backup", 1)
			Expect(err).ToNot(HaveOccurred())
			Expect(content).To(Equal([]byte("version 1 content")))
		})
	})
})
