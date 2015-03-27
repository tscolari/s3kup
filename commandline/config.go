package commandline

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initViperFlags(command *cobra.Command) {
	viper.SetDefault("endpoint-url", "s3.amazonaws.com")
	viper.SetDefault("ssl", true)

	viper.BindPFlag("endpoint-url", command.Flags().Lookup("endpoint-url"))
	viper.BindPFlag("access-key-id", command.Flags().Lookup("access-key-id"))
	viper.BindPFlag("access-key-secret", command.Flags().Lookup("access-key-secret"))
	viper.BindPFlag("bucket-name", command.Flags().Lookup("bucket-name"))
	viper.BindPFlag("file-name", command.Flags().Lookup("file-name"))
	viper.BindPFlag("no-ssl", command.Flags().Lookup("no-ssl"))
	viper.BindPFlag("versions-to-keep", command.Flags().Lookup("versions-to-keep"))
}
