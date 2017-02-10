package main

import (
	"github.com/google/go-github/github"
)
func getLatestDopeRelease() (string, error) {
	client := github.NewClient(nil)
	repoRelease, _, err := client.Repositories.GetLatestRelease("offers", "dope")
	if err != nil {
		return "", err
	}
	return *repoRelease.TagName, nil
}