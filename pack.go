package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/docker/distribution/reference"
	regclient "github.com/docker/distribution/registry/client"
	_ "github.com/motemen/go-loghttp/global"
)

type Pack struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	Cmd     string `json:"cmd"`
	Tag     string `json:"tag"`
	ImageId string `json:"imageId"`
}

func newPack(image string) *Pack {
	parts := strings.Split(image, "/")
	name := parts[len(parts)-1]
	return &Pack{Name: name, Image: image}
}

func repoTags(image string) ([]string, error) {
	ref, err := reference.ParseNamed(image)
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

func highTag(image string) (string, error) {
	tags, err := repoTags(image)
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
func (p *Pack) checkForUpdate() (avail bool, tag string) {
	highTag, err := highTag(p.Image)
	if err != nil {
		log.Error(err)
		return false, ""
	}

	log.Debug("highTag:", highTag)
	if "" == highTag {
		log.Warning("No semantic tags in repo")
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

func (p *Pack) runString() string {
	// {{p.name}}() { {{os.Args[0]}} run {{p.name}} '{{p.cmd}} $@' }
	return fmt.Sprintf("%s() { %s run %s '%s $@' }", p.Name, os.Args[0], p.Name, p.Cmd)
}
