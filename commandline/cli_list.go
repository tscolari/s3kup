package commandline

import (
	"fmt"
	"time"

	"code.cloudfoundry.org/bytefmt"

	"github.com/spf13/cobra"
	"github.com/tscolari/s3kup/list"
	"github.com/tscolari/s3kup/log"
	"github.com/tscolari/s3kup/s3"
)

func listCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List remote stored versions",
		Long:  `List remote stored versions`,
		Run: func(cmd *cobra.Command, args []string) {
			initLogger()
			accessKey, secretKey, bucketName, fileName, endpointURL, err := fetchAndValidateGlobalParams()
			if err != nil {
				log.Fatal(err)
			}

			s3Client := s3.New(accessKey, secretKey, bucketName, endpointURL)
			lister := list.New(s3Client)

			versions, err := lister.List(fileName)
			if err != nil {
				log.Fatal(err)
			}

			if len(versions) == 0 {
				fmt.Println("No versions found")
			}

			for _, version := range versions {
				size := bytefmt.ByteSize(version.Size)
				fmt.Printf("* %d\t%10s\t%s\n", version.Version, size, version.LastModified.Format(time.ANSIC))
			}
		},
	}
	return cmd
}
