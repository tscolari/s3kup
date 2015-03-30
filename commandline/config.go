package commandline

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initViperFlags(mainCmd, pushCmd *cobra.Command) {
	viper.SetDefault("endpoint-url", "https://s3.amazonaws.com")

	viper.BindPFlag("endpoint-url", mainCmd.PersistentFlags().Lookup("endpoint-url"))
	viper.BindPFlag("access-key", mainCmd.PersistentFlags().Lookup("access-key"))
	viper.BindPFlag("secret-key", mainCmd.PersistentFlags().Lookup("secret-key"))
	viper.BindPFlag("bucket-name", mainCmd.PersistentFlags().Lookup("bucket-name"))
	viper.BindPFlag("file-name", mainCmd.PersistentFlags().Lookup("file-name"))
	viper.BindPFlag("verbose", mainCmd.PersistentFlags().Lookup("verbose"))

	viper.BindPFlag("versions-to-keep", pushCmd.Flags().Lookup("versions-to-keep"))
}
