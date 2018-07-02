package cmd

import (
	"github.com/spf13/cobra"

	"fmt"
	"github.com/vamc19/spawner/pkg/image"
	"os"
)

var pullCmd = &cobra.Command{
	Use:   "pull IMAGE",
	Short: "Pull images from registry",
	Args:  cobra.ExactArgs(1),
	Run:   pullImage,
}

func init() {
	rootCmd.AddCommand(pullCmd)
}

func pullImage(cmd *cobra.Command, args []string) {
	i, err := image.InitializeImage(args[0])
	if err != nil {
		os.Exit(1)
	}

	err = i.Pull()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
