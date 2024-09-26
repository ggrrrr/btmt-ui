package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ggrrrr/btmt-ui/be/common/cmd"
	"github.com/ggrrrr/btmt-ui/be/svc-people/cmd/client"
	"github.com/ggrrrr/btmt-ui/be/svc-people/cmd/server"
)

var rootCmd = &cobra.Command{
	Use: "people",
	// Aliases: []string{"insp"},
	Short: "people module",
	// Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var serverCmd = &cobra.Command{
	Use: "server",
	// Aliases: []string{"insp"},
	Short: "start server",
	// Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := server.Server()
		if err != nil {
			fmt.Printf("error %v \n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your '%s'", err)
		os.Exit(1)
	}
}

func main() {
	cmd.Parse(rootCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(client.ClientCmd)
	Execute()
}
