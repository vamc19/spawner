package image

import (
	"errors"
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
}

func InitializeImage(arg string) (Image, error) {
	if arg == "" {
		return Image{}, errors.New("empty argument")
	}

	i := Image{
		User:        defaultUser,
		Tag:         defaultTag,
		RegistryURL: defaultRegistry,
	}

	splitTag := strings.Split(arg, ":")
	if len(splitTag) == 2 {
		i.Tag = splitTag[1]
	}

	splitName := strings.Split(splitTag[0], "/")
	if len(splitName) == 2 {
		i.User = splitName[0]
		i.Repo = splitName[1]
	} else {
		i.Repo = splitName[0]
	}

	return i, nil
}
