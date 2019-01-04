package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func CheckPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { // path exists
		return true, nil
	}
	if os.IsNotExist(err) { // path does not exist
		return false, nil
	}

	return false, err // Some other error, maybe permissions?
}

func ExtractLayer(r io.Reader, path string) error {
	err := os.MkdirAll(path, 644)
	if err != nil {
		return err
	}

	cmd := exec.Command("tar", "-xz", "-C", path)
	cmd.Stdin = r

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Todo: Implement a filesystem interface to support other filesystems
// Layers should be absolute paths ordered with top most layer at index 0
func MountOverlayFS(layers []string, containerDir string, mountPoint string) (string, error) {

	containerDir, _ = filepath.Abs(containerDir)
	upperPath := filepath.Join(containerDir, "upper")
	workPath := filepath.Join(containerDir, "work")

	if mountPoint == "" {
		mountPoint = filepath.Join(containerDir, "mount")
	}

	for _, p := range []string{upperPath, workPath, mountPoint} {
		exists, err := CheckPathExists(p)
		if err != nil {
			return "", err
		}

		if !exists {
			err = os.MkdirAll(p, os.ModePerm)
			if err != nil {
				fmt.Println("Error creating working directory.")
				return "", err
			}
		}
	}

	lowerDirs := strings.Join(layers, ":")
	options := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerDirs, upperPath, workPath)

	err := syscall.Mount("spawner-overlay", mountPoint, "overlay", 0, options)
	if err != nil {
		return "", err
	}

	return mountPoint, nil
}
