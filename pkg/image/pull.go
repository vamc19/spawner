package image

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	authServerURL = "https://auth.docker.io/token"
	authService   = "registry.docker.io"
	manifestType  = "application/vnd.docker.distribution.manifest.v2+json" // media type for v2 image manifest
	httpClient    = &http.Client{Timeout: 5 * time.Second}
)

type Manifest struct { // v2
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"`
	Config        struct {
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
	} `json:"config"`
	Layers []struct {
		MediaType string   `json:"mediaType"`
		Size      int      `json:"size"`
		Digest    string   `json:"digest"`
		Urls      []string `json:"urls"`
	} `json:"layers"`
}

type AuthToken struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
	IssuedAt  string `json:"issued_at"`
}

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

// Check if a manifest for the image already exists on disk
func (i *Image) checkLocalManifest() (bool, error) {
	manifestPath := filepath.Join(i.Store.ManifestPath, i.User, i.Repo, i.Tag+".json")
	_, err := os.Stat(manifestPath)
	if err == nil { // file exists
		return true, nil
	}
	if os.IsNotExist(err) { // file does not exist
		return false, nil
	}

	return false, err // Some other error, maybe permissions?
}

// Pull image from registry
func (i *Image) pullImage() error {
	// Get Token
	t := new(AuthToken)
	err := i.getToken(t)
	if err != nil {
		return err
	}

	// Get Manifest
	m := new(Manifest)
	err = i.getManifest(t.Token, m)
	if err != nil {
		return err
	}

	// Download layers to store
	fmt.Println("Pulling layers...")
	for i, layer := range m.Layers {
		fmt.Printf("%d. %s \t %d\n", i, layer.Digest, layer.Size)
	}

	// write manifest to disk
	err = i.saveManifest(m)
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
	err = ioutil.WriteFile(jsonPath, manifestJson, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Download manifest from registry
func (i *Image) getManifest(token string, m *Manifest) error {
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
func (i *Image) getToken(t *AuthToken) error {
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
