package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type manifest struct {
	packs []pack
}

type pack struct {
	name    string `json:"name"`
	image   string `json:"image"`
	alias   string `json:"alias"`
	version string `json:"version"`
}

func newManifest() *manifest {
	return &manifest{}
}

// Returns the manifest as a JSON byte slice
func (m *manifest) json() ([]byte, error) {
	//TODO test me
	return json.Marshal(m.packs)
}

// Reads from the manifest JSON file on disk
// Returns a manifest
func manifestFromFile(path string) (*manifest, error) {
	//TODO test me
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	packs := []pack{}
	if err = json.Unmarshal(data, &packs); err != nil {
		return nil, err
	}
	return &manifest{packs: packs}, nil
}

func (m *manifest) removePack(name string) {
	//TODO test me
	for i, p := range m.packs {
		if p.name == name {
			m.packs = append(m.packs[:i], m.packs[i+1:]...)
			return
		}
	}
}

func (m *manifest) addPack(p pack) {
	m.packs = append(m.packs, p)
}

func (m *manifest) writeToFile(path string) error {
	if data, err := m.json(); err != nil {
		return ioutil.WriteFile(path, data, 0644)
	} else {
		return err
	}
}

func (m *manifest) writeAliasFile(path string) error {
	//TODO test me
	var buf bytes.Buffer
	for _, p := range m.packs {
		buf.WriteString(fmt.Sprintf("alias %s=%s\n", p.name, p.alias))
	}
	return ioutil.WriteFile(path, buf.Bytes(), 0644)
}

// Check if new version of named package is avilable
// Returns true if so, false otherwise
func (m *manifest) checkForUpdate(name string) (avail bool, image string) {
	p := m.getPack(name)
	if nil != p {
		return p.checkForUpdate()
	}
	return false, ""
}

func (m *manifest) getPack(name string) *pack {
	// TODO test me
	for _, p := range m.packs {
		if p.name == name {
			return &p
		}
	}
	return nil
}

func (p *pack) checkForUpdate() (avail bool, image string) {
	// TODO implement check on p
	// list tags in docker repo
	// find semantically highest tag
	// true if that is higher than our tag
	// return new tag
	return true, "some_tag"
}
