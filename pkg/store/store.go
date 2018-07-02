package store

import (
	"github.com/vamc19/spawner/pkg/utils"
	"os"
	"path/filepath"
)

var (
	folderName     = ".spawner"
	manifestFolder = "manifests"
	layerFolder    = "layers"
)

type Store struct {
	ManifestPath string // path to manifest store
	LayerPath    string // path to layer store
}

func InitializeStore(path string) (Store, error) {
	// Check and create directories
	absPath, err := filepath.Abs(path)
	if err != nil {
		return Store{}, err
	}

	// folders to create
	manifestsPath := filepath.Join(absPath, manifestFolder)
	layersPath := filepath.Join(absPath, layerFolder)

	for _, p := range []string{manifestsPath, layersPath} {
		err = os.MkdirAll(p, os.ModePerm)
		if err != nil {
			return Store{}, err
		}
	}

	return Store{manifestsPath, layersPath}, nil
}

func DefaultStore() (Store, error) {
	userHome, err := utils.GetUserHome()
	if err != nil {
		return Store{}, err
	}

	defaultFolder := filepath.Join(userHome, folderName)
	return InitializeStore(defaultFolder)
}
