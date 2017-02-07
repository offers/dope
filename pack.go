package main

import (
	"fmt"
	"os"

	_ "github.com/chriskite/docker-registry-client/registry"
)

type pack struct {
	name    string `json:"name"`
	image   string `json:"image"`
	cmd     string `json:"cmd"`
	version string `json:"version"`
}

func (p *pack) checkForUpdate() (avail bool, tag string) {
	// TODO implement check on p
	// list tags in docker repo
	// find semantically highest tag
	// true if that is higher than our tag
	// return new tag
	return true, "some_tag"
}

func (p *pack) bashFunction() string {
	// {{p.name}}() { {{os.Args[0]}} run {{p.name}} '{{p.cmd}} $@' }
	return fmt.Sprintf("%s() { %s run %s '%s $@' }", p.name, os.Args[0], p.name, p.cmd)
}
