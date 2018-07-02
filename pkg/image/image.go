package image

import (
	"errors"
	"github.com/vamc19/spawner/pkg/store"
	"strings"
)

var (
	defaultRegistry = "https://registry-1.docker.io/v2"
	defaultTag      = "latest"
	defaultUser     = "library"
)

type Image struct {
	User        string
	Repo        string
	Tag         string
	RegistryURL string
	Store       *store.Store
}

func InitializeImage(arg string) (Image, error) {
	if arg == "" {
		return Image{}, errors.New("empty argument")
	}

	// Default values
	var repo string
	user := defaultUser
	tag := defaultTag
	registry := defaultRegistry // Todo: flag to set registry and store location
	imageStore, err := store.DefaultStore()
	if err != nil {
		return Image{}, err
	}

	// Get image tag if specified
	splitTag := strings.Split(arg, ":")
	if len(splitTag) == 2 {
		tag = splitTag[1]
	}

	// Get username and repo name
	splitName := strings.Split(splitTag[0], "/")
	if len(splitName) == 2 { // arg in the form username/reponame
		user = splitName[0]
		repo = splitName[1]
	} else { // arg only has repo name. username will be set to defaultUser
		repo = splitName[0]
	}

	i := Image{
		User:        user,
		Repo:        repo,
		Tag:         tag,
		RegistryURL: registry,
		Store:       &imageStore,
	}

	return i, nil
}
