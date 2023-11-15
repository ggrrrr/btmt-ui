package cmd

import "github.com/spf13/cobra"

type Flags struct {
	ConfigName  string
	ConfigPaths []string
	LogLevel    string
}

var (
	GlobalFlags Flags
)

func init() {
	GlobalFlags.ConfigName = "config"
	GlobalFlags.ConfigPaths = []string{"./"}
}

func Parse(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().StringVarP(&GlobalFlags.LogLevel, "level", "v", "info", "log level debug/info/error")
	rootCmd.PersistentFlags().StringVarP(&GlobalFlags.ConfigName, "config", "c", "config", "name of the config file without .yaml")
}
