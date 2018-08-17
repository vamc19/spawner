package store

import (
		"os"
	"path/filepath"
)

var (
	defaultFolder   = "/var/lib/spawner"
	manifestFolder  = "manifests"
	layerFolder     = "layers"
	containerFolder = "containers"
)

type Store struct {
	ManifestPath  string // path to manifest store
	LayerPath     string // path to layer store
	ContainerPath string // path to container store
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
	containersPath := filepath.Join(absPath, containerFolder)

	for _, p := range []string{manifestsPath, layersPath, containersPath} {
		err = os.MkdirAll(p, os.ModePerm)
		if err != nil {
			return Store{}, err
		}
	}

	return Store{manifestsPath, layersPath, containersPath}, nil
}

func DefaultStore() (Store, error) {
	return InitializeStore(defaultFolder)
}
