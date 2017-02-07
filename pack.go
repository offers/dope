package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chriskite/docker-registry-client/registry"
)

type pack struct {
	name  string `json:"name"`
	image string `json:"image"`
	cmd   string `json:"cmd"`
	tag   string `json:"tag"`
}

func (p *pack) checkForUpdate() (avail bool, tag string) {
	//TODO don't assume private registry
	parts := strings.Split(p.image, "/")
	repo := parts[0]
	url := fmt.Sprintf("https://%s/", repo)
	username := "" // anonymous
	password := "" // anonymous
	reg, err := registry.New(url, username, password)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}

	tags, err := reg.Tags(fmt.Sprintf("%s/%s", parts[1], parts[2]))
	if err != nil {
		fmt.Println(err)
		return false, ""
	}

	highTag := ""
	for _, t := range tags {
		if compareTags(t, highTag) == 1 {
			highTag = t
		}
	}

	if "" == highTag {
		// no semantic tags in repo
		return false, ""
	}

	if compareTags(highTag, p.tag) == 1 {
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

	for i, _ := range nums1 {
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
	return fmt.Sprintf("%s() { %s run %s '%s $@' }", p.name, os.Args[0], p.name, p.cmd)
}
