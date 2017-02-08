package main

import (
	"os"
	"time"

	"github.com/fsouza/go-dockerclient"
)

func newDockerClient() (*docker.Client, error) {
	endpoint := "unix:///var/run/docker.sock"
	return docker.NewClient(endpoint)
}

func dockerPull(repo string, tag string) error {
	client, err := newDockerClient()
	if err != nil {
		return err
	}

	opts := docker.PullImageOptions{
		Repository:        repo,
		Tag:               tag,
		OutputStream:      os.Stderr,
		InactivityTimeout: time.Duration(30) * time.Second,
	}

	log.Info("Pulling image", repo, ":", tag, "...")
	return client.PullImage(opts, docker.AuthConfiguration{})
}

func dockerRmi(name string) error {
	client, err := newDockerClient()
	if err != nil {
		return err
	}

	return client.RemoveImage(name)
}
