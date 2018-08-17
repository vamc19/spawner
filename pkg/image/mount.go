package image

import (
	"fmt"
	"github.com/vamc19/spawner/pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

func (i *Image) Mount(mountPath string, workPath string) (string, error) {

	var manifest Manifest
	err := i.loadLocalManifest(&manifest)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		// Image does not exist. Pull from registry
		fmt.Println("Pulling image from registry...")

		err = i.Pull()
		if err != nil {
			return "", err
		}
	}

	var layers []string
	for index := len(manifest.Layers) - 1; index >= 0; index-- {
		l := manifest.Layers[index]
		layerName := strings.Split(l.Digest, ":")[1]
		layerPath := filepath.Join(i.Store.LayerPath, layerName)

		layers = append(layers, layerPath)
	}

	mountPath, err = utils.MountOverlayFS(layers, workPath, mountPath)
	if err != nil {
		return "", err
	}

	return mountPath, nil
}
