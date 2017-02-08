package main

import (
	"os"
	"time"

	"github.com/fsouza/go-dockerclient"
)

func dockerPull(repo string, tag string) error {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
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
