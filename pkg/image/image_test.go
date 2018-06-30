package image

import (
	"github.com/vamc19/spawner/pkg/utils"
	"testing"
)

func TestInitializeImage(t *testing.T) {
	i, err := InitializeImage("vamc19/myfirstapp:latest")
	utils.AssertEqual(t, err, nil, "Should not be nil when non empty string is passed")
	utils.AssertEqual(t, i.Repo, "myfirstapp", "Wrong repository")
	utils.AssertEqual(t, i.User, "vamc19", "Wrong Docker Hub User")
	utils.AssertEqual(t, i.Tag, "latest", "Wrong Tag")

	i, err = InitializeImage("ubuntu:16.04")
	utils.AssertEqual(t, err, nil, "Should not be nil when non empty string is passed")
	utils.AssertEqual(t, i.Repo, "ubuntu", "Wrong repository")
	utils.AssertEqual(t, i.User, "library", "Wrong Docker Hub User")
	utils.AssertEqual(t, i.Tag, "16.04", "Wrong Tag")

	i, err = InitializeImage("")
	utils.AssertNotEqual(t, err, nil, "Should be nil when empty string is passed")
	utils.AssertEqual(t, i, Image{}, "Should be empty")
}
