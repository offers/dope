package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"encoding/json"
	"github.com/mitchellh/go-homedir"
	"github.com/offers/dope/out"
	"github.com/op/go-logging"
	"gopkg.in/urfave/cli.v1"
	"os/exec"
)

const Release = "0.0.7"

var log = logging.MustGetLogger("dope")

func initConfDir() string {
	home, err := homedir.Dir()
	if err != nil {
		out.Error(err)
		os.Exit(1)
	}
	confDir := filepath.Join(home, ".dope")
	os.MkdirAll(confDir, 0755)
	return confDir
}

func setupLogging() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
}

func main() {
	app := cli.NewApp()
	app.Name = "dope"
	app.Usage = "Simple docker package management"
	app.Version = Release

	setupLogging()

	confDir := initConfDir()
	manifest, err := initManifest(confDir)
	if err != nil {
		out.Error(err)
		os.Exit(1)
	}

	app.Commands = []cli.Command{
		{
			Name:    "self-update",
			Aliases: []string{"sup"},
			Usage:   "update dope",
			Action: func(c *cli.Context) error {
				avail, err := selfUpdateAvail()
				if err != nil {
					out.Error(err)
					out.Notice("Couldn't get the latest release for dope")
					return nil
				}

				if !avail {
					out.Info("dope is already up-to-date")
					return nil
				}

				out.Info("Installing latest version of dope...")
				output, err := exec.Command("bash", "-c", "\\curl -sSL https://raw.githubusercontent.com/offers/dope/master/install.sh | sudo bash").CombinedOutput()
				if err != nil {
					out.Error(err)
				}
				fmt.Printf("%s", output)
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list installed packages",
			Action: func(c *cli.Context) error {
				notifyIfSelfUpdateAvail()

				out.Info("Installed packages:\n")
				for _, p := range manifest.Packs {
					out.Println(p.Name, p.Tag)
				}
				return nil
			},
		},
		{
			Name:    "update",
			Aliases: []string{"up"},
			Usage:   "update packages",
			Action: func(c *cli.Context) error {
				notifyIfSelfUpdateAvail()

				if c.NArg() > 0 {
					// update single package
					name := c.Args()[0]
					p := manifest.getPack(name)
					if nil == p {
						out.Notice("No package named", name, "is installed")
						return err
					}

					return updatePack(manifest, p)
				} else {
					// update all packages
					updateAllPacks(manifest)
				}
				return nil
			},
		},
		{
			Name:    "install",
			Aliases: []string{"in"},
			Usage:   "install a package",
			Action: func(c *cli.Context) (err error) {
				notifyIfSelfUpdateAvail()

				if c.NArg() > 0 {
					repo := c.Args()[0]

					if manifest.isInstalled(repo) {
						out.Notice(repo, "already installed, try update instead")
						return nil
					}

					// install package
					pack, err := installImage(repo)
					if err != nil {
						out.Error(err)
						return err
					}

					if err := manifest.addPack(pack); err != nil {
						out.Error(err)
						return err
					}

					out.Success("Installed", pack.Name, pack.Tag)
				} else {
					err = errors.New("No package name given to install")
					out.Error(err)
					return err
				}
				return nil
			},
		},
		{
			Name:    "uninstall",
			Aliases: []string{"un"},
			Usage:   "uninstall a package",
			Action: func(c *cli.Context) (err error) {
				notifyIfSelfUpdateAvail()

				if c.NArg() > 0 {
					name := c.Args()[0]

					pack, err := manifest.removePackWithName(name)
					if err != nil {
						out.Error(err)
						return nil
					}

					if nil == pack {
						out.Notice(name, "is not installed")
						return nil
					}

					// rm docker image
					err = removeImage(pack.Repo, pack.Tag)
					if err != nil {
						return err
					}

					out.Success("Uninstalled", pack.Name)
				} else {
					err = errors.New("no package name given to install")
					out.Error(err)
					return nil
				}
				return nil
			},
		},
		{
			Name:    "check",
			Aliases: []string{"ch"},
			Usage:   "check for updates to package",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "q, quiet",
					Usage: "only output if an update is available",
				},
			},
			SkipFlagParsing: false,
			Action: func(c *cli.Context) (err error) {
				notifyIfSelfUpdateAvail()

				if c.NArg() > 0 {
					// check package for updates
					name := c.Args()[0]
					avail, _, tag := manifest.checkForUpdate(name)
					if avail {
						out.Info("New version", tag, "available for", name)
					} else if !c.Bool("quiet") {
						out.Info("No update available for", name)
					}
				} else {
					out.Notice("no package name given to check")
				}

				return nil
			},
		},
	}

	app.Run(os.Args)
}

func notifyIfSelfUpdateAvail() {
	avail, _ := selfUpdateAvail()
	if avail {
		out.Info("A new version of dope is available. Run 'dope self-update' to install.")
	}
}

func selfUpdateAvail() (bool, error) {
	tag, err := getLatestDopeRelease()
	if err != nil {
		return false, err
	}

	if 1 == compareTags(tag, Release) {
		return true, nil
	}
	return false, nil
}

func installImage(repo string) (*Pack, error) {
	tag, err := highTag(repo)
	if err != nil {
		return nil, err
	}

	err = dockerPull(repo, tag)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(repo, "/")
	if len(parts) < 1 {
		return nil, errors.New("invalid repo: " + repo)
	}
	name := parts[len(parts)-1]

	dopejson, err := dockerGetDopeFile(repo, tag)
	if err != nil {
		// no .dope.json found
		out.Notice("No .dope.json found, using dumb defaults") // TODO get info from user
	} else {
		pack := &Pack{Name: name}
		err = json.Unmarshal(dopejson, pack)
		if err != nil {
			return nil, err
		}
		pack.Tag = tag
		pack.Repo = repo
		return pack, nil
	}
	return newDefaultPack(repo, tag, name)
}

func removeImage(repo string, tag string) error {
	image := fmt.Sprintf("%s:%s", repo, tag)
	out.Println("Removing Docker image", image, "...")
	return dockerRmi(image)
}

func updatePack(m *Manifest, pack *Pack) error {
	avail, repo, tag := m.checkForUpdate(pack.Name)
	if avail {
		out.Info("New version", tag, "available for", pack.Name)
		newPack, err := installImage(repo)
		if err != nil {
			out.Error(err)
			return err
		}

		oldPack, err := m.removePack(pack)
		if err != nil {
			out.Error(err)
			return err
		}
		removeImage(repo, oldPack.Tag)

		m.addPack(newPack)
		out.Successf("Updated %s from %s to %s\n", newPack.Name, oldPack.Tag, newPack.Tag)
	} else {
		out.Info("No update available for", pack.Name)
	}

	return nil
}

// TODO return []error
func updateAllPacks(m *Manifest) {
	packs := m.Packs
	for _, p := range packs {
		updatePack(m, p)
	}
}
