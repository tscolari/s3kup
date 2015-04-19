package commandline

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tscolari/s3kup/backup"
	"github.com/tscolari/s3kup/log"
	"github.com/tscolari/s3kup/s3"
)

func pushCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Pushes the piped input to s3",
		Long:  `Pushes the pipped input to s3, as a versioned backup`,
		Run: func(cmd *cobra.Command, args []string) {
			initLogger()
			accessKey, secretKey, bucketName, fileName, endpointURL, err := fetchAndValidateGlobalParams()
			if err != nil {
				log.Fatal(err)
			}
			versionsToKeep, err := fetchVersionsToKeep()
			if err != nil {
				log.Fatal(err)
			}

			s3Client := s3.New(accessKey, secretKey, bucketName, endpointURL)
			backuper := backup.New(s3Client, versionsToKeep)

			content, err := getInputContent()
			if err != nil {
				log.Fatal(err)
			}

			err = backuper.Backup(fileName, content)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	cmd.Flags().IntP("versions-to-keep", "k", 5, "Number of versions to keep")
	return cmd
}

func getInputContent() ([]byte, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if fi.Mode()&os.ModeNamedPipe != 0 {
		return ioutil.ReadAll(os.Stdin)
	}

	return nil, errors.New("not using pipeline")
}

func fetchVersionsToKeep() (versionsToKeep int, err error) {
	if versionsToKeep = viper.GetInt("versions-to-keep"); versionsToKeep <= 0 {
		err = errors.New("invalid versions to keep. Must be 1 or greater")
	}

	return versionsToKeep, err
}
