package image

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vamc19/spawner/pkg/utils"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type Manifest struct { // v2
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"`
	Config        struct {
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
	} `json:"config"`
	Layers []layer `json:"layers"`
}

type layer struct {
	MediaType string   `json:"mediaType"`
	Size      int      `json:"size"`
	Digest    string   `json:"digest"`
	Urls      []string `json:"urls"`
}

type authToken struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
	IssuedAt  string `json:"issued_at"`
}

// Check if a manifest for the image already exists on disk
func (i *Image) checkLocalManifest() (bool, error) {
	return utils.CheckPathExists(i.getManifestPath())
}

// Load manifest from store
func (i *Image) loadLocalManifest(m *Manifest) error {
	manifestPath := i.getManifestPath()
	manifestFile, err := os.Open(manifestPath)
	if err != nil {
		return err
	}

	defer manifestFile.Close()

	bv, err := ioutil.ReadAll(manifestFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bv, &m)
	if err != nil {
		return err
	}

	return nil
}

// Todo: rkt implements CAS for manifests. interesting. See how to do that.
// Manifest will be saved in ~/.spawner/manifests/<registry username>/<registry repo>/<tagname>.json
func (i *Image) saveManifest(m *Manifest) error {
	manifestJson, err := json.Marshal(m)
	if err != nil {
		return err
	}

	jsonFolder := filepath.Join(i.Store.ManifestPath, i.User, i.Repo)
	err = os.MkdirAll(jsonFolder, os.ModePerm)
	if err != nil {
		return err
	}

	jsonPath := filepath.Join(jsonFolder, i.Tag+".json")
	err = ioutil.WriteFile(jsonPath, manifestJson, 0666)
	if err != nil {
		return err
	}

	return nil
}

// Download manifest from registry
func (i *Image) downloadManifest(token string, m *Manifest) error {
	// example url: https://registry-1.docker.io/v2/library/ubuntu/manifests/latest
	url := fmt.Sprintf("%s/%s/%s/manifests/%s", i.RegistryURL, i.User, i.Repo, i.Tag)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", manifestType)
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	contentType := res.Header["Content-Type"][0]
	if contentType != manifestType {
		return errors.New("only v2 images are supported\n")
	}

	err = json.NewDecoder(res.Body).Decode(m)
	if err != nil {
		return err
	}

	return nil
}

// Get bearer token from auth server
func (i *Image) getToken(t *authToken) error {
	// example url: https://auth.docker.io/token?service=registry.docker.io&scope=repository:library/ubuntu:pull
	url := fmt.Sprintf("%s?service=%s&scope=repository:%s/%s:pull", authServerURL, authService, i.User, i.Repo)
	res, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(t)
	if err != nil {
		return err
	}

	return nil
}

func (i *Image) getManifestPath() string {
	return filepath.Join(i.Store.ManifestPath, i.User, i.Repo, i.Tag+".json")
}