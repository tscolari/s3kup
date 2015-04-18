package commandline

import (
	"encoding/binary"
	"os"
	"strconv"

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

			if len(args) > 1 {
				log.Fatal("You can specify only one version to get")
			}

			var content []byte
			if len(args) == 1 {
				var version int64
				version, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					log.Fatal("Invalid version format. It can only contain numbers")
				}

				content, err = fetcher.FetchVersion(fileName, version)
			} else {
				content, err = fetcher.FetchLatest(fileName)
			}

			if err != nil {
				log.Fatal(err)
			}

			binary.Write(os.Stdout, binary.LittleEndian, content)
		},
	}
	return cmd
}
