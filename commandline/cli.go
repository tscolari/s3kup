package commandline

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tscolari/s3up/backup"
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
			accessKeyID := viper.GetString("access-key-id")
			accessKeySecret := viper.GetString("access-key-secret")
			bucketName := viper.GetString("bucket-name")
			versionsToKeep := viper.GetInt("versions-to-keep")
			fileName := viper.GetString("file-name")
			endpointURL := getEndpointURL()

			s3Client := s3.New(accessKeyID, accessKeySecret, bucketName, endpointURL)
			backuper := backup.New(s3Client, versionsToKeep)

			content, err := getInputContent()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = backuper.Backup(fileName, content)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
		},
	}

	cmd.Flags().StringP("endpoint-url", "e", "s3.amazonaws.com", "the s3 region endpoint url (see http://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region)")
	cmd.Flags().StringP("access-key-id", "i", "", "AWS Access Key ID")
	cmd.Flags().StringP("access-key-secret", "s", "", "AWS Access Key Secret")
	cmd.Flags().StringP("bucket-name", "b", "", "Target S3 bucket")
	cmd.Flags().StringP("file-name", "n", "", "How the file will be called on s3")
	cmd.Flags().Bool("no-ssl", false, "Use ssl endpoint")
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

func getEndpointURL() string {
	endpointURL := viper.GetString("endpoint-url")
	ssl := !viper.GetBool("no-ssl")

	if ssl {
		return fmt.Sprintf("https://%s", endpointURL)
	}
	return fmt.Sprintf("http://%s", endpointURL)
}
