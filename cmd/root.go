package cmd

import (
	"os"

	"errors"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "spawner [command]",
	Short: "Run docker images in systemd's namespace containers.",
	Long:  `A quick and dirty way of running Docker images using systemd-nspawn`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("spawner: missing command")
	},
}

// Execute adds all child cmd to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
