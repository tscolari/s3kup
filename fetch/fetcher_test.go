package fetch_test

import (
	"errors"

	"github.com/tscolari/s3up/fetch"
	"github.com/tscolari/s3up/s3"
	"github.com/tscolari/s3up/s3/fakeclient"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fetcher", func() {
	var fetcher fetch.Fetcher
	var client *fakeclient.Client

	BeforeEach(func() {
		client = &fakeclient.Client{}
		fetcher = fetch.New(client)
	})

	Describe("fetch", func() {
		BeforeEach(func() {
			client.ListCall = func(path string) (s3.Versions, error) {
				return s3.Versions{
					s3.Version{Path: "my-backup/0", Version: 0},
					s3.Version{Path: "my-backup/1", Version: 1},
					s3.Version{Path: "my-backup/2", Version: 2},
				}, nil
			}
		})

		Context("when the s3 client returns an error", func() {
			It("forwards the error", func() {
				client.GetCall = func(path string) ([]byte, error) {
					return nil, errors.New("some error")
				}

				_, err := fetcher.Fetch("my-backup")
				Expect(err).To(MatchError("some error"))
			})
		})

		Context("when there is no versions/backup", func() {
			It("returns an error", func() {
				client.ListCall = func(path string) (s3.Versions, error) {
					return s3.Versions{}, nil
				}
				_, err := fetcher.Fetch("dontexist")
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

			content, err := fetcher.Fetch("my-backup")
			Expect(err).ToNot(HaveOccurred())
			Expect(content).To(Equal([]byte("correct version")))
		})
	})
})
