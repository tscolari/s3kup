package commandline

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tscolari/s3up/fetch"
	"github.com/tscolari/s3up/log"
	"github.com/tscolari/s3up/s3"
)

func pullCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Get remote version contents",
		Long:  `Get remote version and print it's contents to STDOUT`,
		Run: func(cmd *cobra.Command, args []string) {
			initLogger()
			accessKey, secretKey, bucketName, fileName, endpointURL, err := fetchAndValidateGlobalParams()
			if err != nil {
				log.Fatal(err)
			}

			s3Client := s3.New(accessKey, secretKey, bucketName, endpointURL)
			fetcher := fetch.New(s3Client)

			content, err := fetcher.FetchLatest(fileName)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf(string(content))
		},
	}
	return cmd
}
