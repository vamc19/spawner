package cmd

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"github.com/vamc19/spawner/pkg/image"
	"os"
	"path/filepath"
)

var (
	mountCmd = &cobra.Command{
		Use:   "mount IMAGE [PATH]",
		Short: "Mount image layers",
		Args:  cobra.MinimumNArgs(1),
		Run:   mountImage,
	}

	workDir string
)

func init() {
	mountCmd.Flags().StringVarP(&workDir, "work-dir", "w", "", "specify working directory")
	rootCmd.AddCommand(mountCmd)
}

func mountImage(cmd *cobra.Command, args []string) {
	i, err := image.InitializeImage(args[0])
	if err != nil {
		os.Exit(1)
	}

	mountPath := ""
	if len(args) > 1 {
		mountPath = args[1]
	}

	if workDir == "" {
		tmpId := uuid.NewV4()
		workDir = filepath.Join(os.TempDir(), tmpId.String())
	}

	mountPath, err = i.Mount(mountPath, workDir)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Println("File system available at: ", mountPath)
}
