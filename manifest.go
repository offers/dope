package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

type manifest struct {
	version string `json:"version"`
	packs   []pack `json:"packs"`
}

func newManifest() *manifest {
	return &manifest{}
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
	if data, err := json.Marshal(m); err != nil {
		return ioutil.WriteFile(path, data, 0644)
	} else {
		return err
	}
}

func (m *manifest) writeAliasFile(path string) error {
	// TODO write bash functions to file
	// to prevent cli from interpreting flags in dope run
	var buf bytes.Buffer
	for _, p := range m.packs {
		buf.WriteString(p.bashFunction() + "\n")
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
