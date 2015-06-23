package commandline

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tscolari/s3kup/log"
)

func New() *cobra.Command {
	mainCmd := mainCommand()
	pushCmd := pushCommand()
	listCmd := listCommand()
	pullCmd := pullCommand()

	mainCmd.AddCommand(pushCmd)
	mainCmd.AddCommand(listCmd)
	mainCmd.AddCommand(pullCmd)

	setGlobalFlags(mainCmd)
	initViperFlags(mainCmd, pushCmd)
	return mainCmd
}

func mainCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "s3kup",
		Run: func(cmd *cobra.Command, args []string) {
			initLogger()
			_, _, _, _, _, err := fetchAndValidateGlobalParams()
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return cmd
}

func setGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("endpoint-url", "e", "https://s3.amazonaws.com", "the s3 region endpoint url (see http://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region)")
	cmd.PersistentFlags().StringP("access-key", "a", "", "AWS Access Key")
	cmd.PersistentFlags().StringP("secret-key", "s", "", "AWS Secret Key")
	cmd.PersistentFlags().StringP("bucket-name", "b", "", "Target S3 bucket")
	cmd.PersistentFlags().StringP("file-name", "n", "", "How the file will be called on s3")
	cmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose mode")
}

func fetchAndValidateGlobalParams() (accessKey, secretKey, bucketName, fileName, endpointURL string, err error) {
	if accessKey = viper.GetString("access-key"); accessKey == "" {
		err = errors.New("missing access key argument")
	}

	if secretKey = viper.GetString("secret-key"); secretKey == "" {
		err = errors.New("missing secret key argument")
	}

	if fileName = viper.GetString("file-name"); fileName == "" {
		err = errors.New("missing file name argument")
	}

	if bucketName = viper.GetString("bucket-name"); bucketName == "" {
		err = errors.New("missing bucket name argument")
	}

	endpointURL = viper.GetString("endpoint-url")

	return accessKey, secretKey, bucketName, fileName, endpointURL, err
}

func initLogger() {
	if viper.GetBool("verbose") {
		log.SetLevel(log.INFO_LEVEL)
	}
}
