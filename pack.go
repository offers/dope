package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/docker/distribution/reference"
	regclient "github.com/docker/distribution/registry/client"
	"github.com/mitchellh/go-homedir"
	"github.com/offers/dope/out"
)

type Pack struct {
	Name       string   `json:"name"`
	Repo       string   `json:"repo"`
	DockerArgs string `json:"dockerArgs"`
	Cmd        string   `json:"cmd"`
	Tag        string   `json:"tag"`
	ImageId    string   `json:"ImageId"`
}

func newDefaultPack(repo string, tag string, name string) (*Pack, error) {
	p := &Pack{
		Name:       name,
		Repo:       repo,
		Tag:        tag,
		DockerArgs: "-it",
	}
	return p, nil
}

func repoTags(repo string) ([]string, error) {
	ref, err := reference.ParseNamed(repo)
	if err != nil {
		return []string{}, err
	}

	url := fmt.Sprintf("https://%s/", reference.Domain(ref))
	name, err := reference.WithName(reference.Path(ref))
	if err != nil {
		return []string{}, err
	}

	r, err := regclient.NewRepository(nil, name, url, http.DefaultTransport)
	if err != nil {
		return []string{}, err
	}

	return r.Tags(nil).All(nil)
}

func highTag(repo string) (string, error) {
	tags, err := repoTags(repo)
	if err != nil {
		return "", err
	}

	highTag := ""
	for _, t := range tags {
		match, _ := regexp.MatchString(`\d\.\d\.\d`, t)
		if match && (highTag == "" || compareTags(t, highTag) == 1) {
			highTag = t
		}
	}

	return highTag, nil
}

// TODO handle versions starting with v, e.g. v1.0.0
// TODO use 'latest' if no semantic tags, and check image hash for update
// TODO add error return
func (p *Pack) checkForUpdate() (avail bool, tag string) {
	highTag, err := highTag(p.Repo)
	if err != nil {
		return false, ""
	}

	if "" == highTag {
		out.Notice("No semantic tags in repo")
		return false, ""
	}

	if compareTags(highTag, p.Tag) == 1 {
		return true, highTag
	}

	return false, ""
}

// Compare 2 semantic version tags (e.g. 1.0.0)
// Returns 1 if t1 is a higher version
// Returns -1 if t2 is a higher version
// Returns 0 if t1 == t2
func compareTags(t1 string, t2 string) int {
	nums1 := strings.Split(t1, ".")
	nums2 := strings.Split(t2, ".")

	l := len(nums1)
	if len(nums2) < l {
		l = len(nums2)
	}

	for i := 0; i < l; i++ {
		v1, _ := strconv.ParseInt(nums1[i], 0, 64)
		v2, _ := strconv.ParseInt(nums2[i], 0, 64)

		if v1 > v2 {
			return 1
		}
		if v1 < v2 {
			return -1
		}
	}

	return 0
}

func (p *Pack) stubPath() string {
	home, err := homedir.Dir()
	if err != nil {
		panic("Couldn't determine user home directory")
	}
	binDir := filepath.Join(home, ".dope", "bin")
	os.MkdirAll(binDir, 0755)
	return filepath.Join(binDir, p.Name)
}

func (p *Pack) removeStub() error {
	return os.Remove(p.stubPath())
}

func (p *Pack) writeStub() error {
	dockerBin, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	dockerCmd := fmt.Sprintf("%s run --rm %s %s:%s", dockerBin, p.DockerArgs, p.Repo, p.Tag)
	if "" != p.Cmd {
		dockerCmd += " " + p.Cmd
	}

	s := fmt.Sprintf("#!/bin/bash\ndope check -q %s\nexec %s $@", p.Name, dockerCmd)
	return ioutil.WriteFile(p.stubPath(), []byte(s), 0755)
}
