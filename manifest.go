package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

const ManifestVersion int64 = 1

type Manifest struct {
	filename string
	Version  int64   `json:"version"`
	Packs    []*Pack `json:"packs"`
}

func initManifest(confDir string) (*Manifest, error) {
	manifestPath := filepath.Join(confDir, "manifest.json")

	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		m := newManifest(manifestPath)
		if err := m.writeToFile(); err != nil {
			return nil, err
		}
		return m, nil
	} else {
		return manifestFromFile(manifestPath)
	}
}

func newManifest(path string) *Manifest {
	return &Manifest{Version: ManifestVersion, filename: path}
}

// Reads from the manifest JSON file on disk
// Returns a manifest
func manifestFromFile(path string) (*Manifest, error) {
	//TODO test me
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	manifest := Manifest{}
	if err = json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	manifest.filename = path
	return &manifest, nil
}

// Removes pack from the manifest and deletes its stub
func (m *Manifest) removePack(pack *Pack) (*Pack, error) {
	return m.removePackWithName(pack.Name)
}

// Removes pack from the manifest by name and deletes its stub
func (m *Manifest) removePackWithName(name string) (*Pack, error) {
	//TODO test me
	for i, p := range m.Packs {
		if p.Name == name {
			m.Packs = append(m.Packs[:i], m.Packs[i+1:]...)
			err := m.writeToFile()
			p.removeStub()
			return p, err
		}
	}
	return nil, nil
}

func (m *Manifest) addPack(p *Pack) error {
	m.Packs = append(m.Packs, p)
	if err := p.writeStub(); err != nil {
		return err
	}
	return m.writeToFile()
}

func (m *Manifest) writeToFile() error {
	if "" == m.filename {
		return errors.New("no filename set for manifest")
	}

	if data, err := json.Marshal(m); err != nil {
		return err
	} else {
		return ioutil.WriteFile(m.filename, data, 0644)
	}
}

// For each pack, write a file to /usr/local/bin which
// will execute that docker repo through dope
func (m *Manifest) writeAliasFiles() error {
	//TODO implement me
	return nil
}

// Check if new version of named package is available
// Returns true if so, false otherwise
func (m *Manifest) checkForUpdate(name string) (avail bool, repo string, tag string) {
	p := m.getPack(name)
	if nil != p {

		avail, tag := p.checkForUpdate()
		return avail, p.Repo, tag
	}
	return false, "", ""
}

func (m *Manifest) getPack(name string) *Pack {
	// TODO test me
	for _, p := range m.Packs {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func (m *Manifest) isInstalled(repo string) bool {
	for _, p := range m.Packs {
		if p.Repo == repo {
			return true
		}
	}
	return false
}
