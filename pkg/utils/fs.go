package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"os/user"
	"path/filepath"
)

func GetUserHome() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir, nil
}

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
