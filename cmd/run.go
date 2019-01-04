package cmd

import (
	"github.com/spf13/cobra"

	"fmt"
	"github.com/vamc19/spawner/pkg/image"
	"os"
)

var runCmd = &cobra.Command{
	Use:   "run IMAGE CMD",
	Short: "Start a container",
	Args:  cobra.ExactArgs(2),
	Run:   startContainer,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func startContainer(cmd *cobra.Command, args []string) {
	i, err := image.InitializeImage(args[0])
	if err != nil {
		os.Exit(1)
	}

	//err = i.Run(args[1])
	//if err != nil {
	//	fmt.Print(err)
	//	os.Exit(1)
	//}
	fmt.Print(i.RegistryURL)
}
