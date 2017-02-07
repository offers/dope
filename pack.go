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

type pack struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	Cmd   string `json:"cmd"`
	Tag   string `json:"tag"`
}

func newPack(image string) *pack {
	parts := strings.Split(image, "/")
	name := parts[len(parts)-1]
	return &pack{Name: name, Image: image}
}

func getRepoTags(image string) ([]string, error) {
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

func (p *pack) checkForUpdate() (avail bool, tag string) {
	tags, err := getRepoTags(p.Image)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}

	highTag := ""
	for _, t := range tags {
		match, _ := regexp.MatchString(`\d\.\d\.\d`, t)
		if match && compareTags(t, highTag) == 1 {
			highTag = t
		}
	}

	if "" == highTag {
		fmt.Println("no semantic tags in repo")
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

func (p *pack) bashFunction() string {
	// {{p.name}}() { {{os.Args[0]}} run {{p.name}} '{{p.cmd}} $@' }
	return fmt.Sprintf("%s() { %s run %s '%s $@' }", p.Name, os.Args[0], p.Name, p.Cmd)
}
