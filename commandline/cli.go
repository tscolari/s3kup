package commandline

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tscolari/s3up/backup"
	"github.com/tscolari/s3up/log"
	"github.com/tscolari/s3up/s3"
)

func New() *cobra.Command {
	mainCmd := mainCommand()

	initViperFlags(mainCmd)
	return mainCmd
}

func mainCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "s3up",
		Short: "simple single file s3 backup tool",
		Long:  `An easy way to backup any file/command output to a s3 bucket.`,
		Run: func(cmd *cobra.Command, args []string) {
			initLogger()
			accessKey, secretKey, bucketName, fileName, endpointURL, versionsToKeep, err := fetchAndValidateParams()
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

	cmd.Flags().StringP("endpoint-url", "e", "https://s3.amazonaws.com", "the s3 region endpoint url (see http://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region)")
	cmd.Flags().StringP("access-key", "i", "", "AWS Access Key")
	cmd.Flags().StringP("secret-key", "s", "", "AWS Secret Key")
	cmd.Flags().StringP("bucket-name", "b", "", "Target S3 bucket")
	cmd.Flags().StringP("file-name", "n", "", "How the file will be called on s3")
	cmd.Flags().IntP("versions-to-keep", "k", 5, "Number of versions to keep")
	cmd.Flags().BoolP("verbose", "v", false, "Verbose mode")

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

func fetchAndValidateParams() (accessKey, secretKey, bucketName, fileName, endpointURL string, versionsToKeep int, err error) {
	if accessKey = viper.GetString("access-key"); accessKey == "" {
		err = errors.New("missing access key argument")
	}

	if secretKey = viper.GetString("secret-key"); secretKey == "" {
		err = errors.New("missing secret key argument")
	}

	if versionsToKeep = viper.GetInt("versions-to-keep"); versionsToKeep <= 0 {
		err = errors.New("invalid versions to keep. Must be 1 or greater")
	}

	if fileName = viper.GetString("file-name"); fileName == "" {
		err = errors.New("missing file name argument")
	}

	if bucketName = viper.GetString("bucket-name"); bucketName == "" {
		err = errors.New("missing bucket name argument")
	}

	endpointURL = viper.GetString("endpoint-url")

	return accessKey, secretKey, bucketName, fileName, endpointURL, versionsToKeep, err
}

func initLogger() {
	if viper.GetBool("verbose") {
		log.SetLevel(log.INFO_LEVEL)
	}
}
