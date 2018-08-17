package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
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
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzipReader)
	for { // iterate over the contents of archive
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// path of the current item
		dst := filepath.Join(path, header.Name)

		if header.Typeflag == tar.TypeDir {
			err = os.MkdirAll(dst, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			continue
		}

		// If it's not a directory, it has to be a file (?)
		f, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, tarReader)
		if err != nil {
			return err
		}
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
			err = os.MkdirAll(p, 0755)
			if err != nil {
				fmt.Println("Error creating working directory.")
				return "", err
			}
		}
	}

	lowerDirs := strings.Join(layers, ":")
	options := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerDirs, upperPath, workPath)

	err := syscall.Mount("none", mountPoint, "overlay", 0, options)
	if err != nil {
		return "", err
	}

	return mountPoint, nil
}
