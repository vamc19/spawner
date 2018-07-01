package image

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
	// Todo: Check if the image is in imagestore already

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
	// Todo: write manifest file to disk

	// Download layers to imagestore
	fmt.Println("Pulling layers...")
	for i, layer := range m.Layers {
		fmt.Printf("%d. %s \t %d\n", i, layer.Digest, layer.Size)
	}

	return nil
}

func (i *Image) getManifest(token string, m *Manifest) error {

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
