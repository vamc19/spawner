package image

import (
	"fmt"
	"github.com/vamc19/spawner/pkg/utils"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	authServerURL = "https://auth.docker.io/token"
	authService   = "registry.docker.io"
	manifestType  = "application/vnd.docker.distribution.manifest.v2+json" // media type for v2 image manifest
	httpClient    = &http.Client{}
)

func (i *Image) Pull() error {
	// Check if the image is in store already
	exists, err := i.checkLocalManifest()
	if err != nil {
		return err
	}

	// Todo: update policy?
	if exists {
		fmt.Println("Image exists")
		return nil
	}

	err = i.pullImage()
	if err != nil {
		return err
	}

	return nil
}

// Pull image from registry
func (i *Image) pullImage() error {
	// Get Token
	var t authToken
	err := i.getToken(&t)
	if err != nil {
		return err
	}

	// Get Manifest
	var m Manifest
	err = i.downloadManifest(t.Token, &m)
	if err != nil {
		return err
	}

	// Download layers to store
	fmt.Println("Pulling layers...")
	for _, l := range m.Layers {
		err = i.pullLayer(l, t.Token)
		if err != nil {
			return err
		}
	}

	// write manifest to disk
	err = i.saveManifest(&m)
	if err != nil {
		return err
	}

	return nil
}

// Download layers.
func (i *Image) pullLayer(l layer, token string) error {
	fmt.Printf("Downloading %s... \t", l.Digest)
	layerName := strings.Split(l.Digest, ":")[1]
	layerPath := filepath.Join(i.Store.LayerPath, layerName)
	exists, err := utils.CheckPathExists(layerPath)
	if err != nil {
		return err
	}

	if exists {
		fmt.Printf("Layer already downloaded\n")
		return nil
	}

	// Download tar
	url := fmt.Sprintf("%s/%s/%s/blobs/%s", i.RegistryURL, i.User, i.Repo, l.Digest)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", l.MediaType)
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = utils.ExtractLayer(res.Body, layerPath)
	if err != nil {
		return err
	}
	fmt.Printf("Done\n")

	return nil
}
