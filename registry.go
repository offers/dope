package main

type registry struct {
	packs []pack
}

type pack struct {
	name    string
	url     string
	image   string
	alias   string
	version string
}

func newRegistry() *registry {
	return &registry{}
}

func registryFromFile(path string) *registry {
	// TODO implement me
	return &registry{}
}

func (r *registry) removePack(name string) {
	//TODO test me
	for i, p := range r.packs {
		if p.name == name {
			r.packs = append(r.packs[:i], r.packs[i+1:]...)
			return
		}
	}
}

func (r *registry) addPack(p pack) {
	r.packs = append(r.packs, p)
}

func (r *registry) writeToFile(path string) error {
	//TODO implement me
	return nil
}

func (r *registry) writeAliasFile(path string) error {
	//TODO implement me
	return nil
}

// Check if new version of named package is avilable
// Returns true if so, false otherwise
func (r *registry) checkForUpdate(name string) (avail bool, image string) {
	p := r.getPack(name)
	if nil != p {
		return p.checkForUpdate()
	}
	return false, ""
}

func (r *registry) getPack(name string) *pack {
	// TODO test me
	for _, p := range r.packs {
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
