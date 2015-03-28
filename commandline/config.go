package commandline

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initViperFlags(command *cobra.Command) {
	viper.SetDefault("endpoint-url", "s3.amazonaws.com")
	viper.SetDefault("ssl", true)

	viper.BindPFlag("endpoint-url", command.Flags().Lookup("endpoint-url"))
	viper.BindPFlag("access-key", command.Flags().Lookup("access-key"))
	viper.BindPFlag("secret-key", command.Flags().Lookup("secret-key"))
	viper.BindPFlag("bucket-name", command.Flags().Lookup("bucket-name"))
	viper.BindPFlag("file-name", command.Flags().Lookup("file-name"))
	viper.BindPFlag("versions-to-keep", command.Flags().Lookup("versions-to-keep"))
}
