package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/op/go-logging"
	"gopkg.in/urfave/cli.v1"
)

var log = logging.MustGetLogger("dope")

func initConfDir() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	confDir := filepath.Join(home, ".dope")
	os.MkdirAll(confDir, 0755)
	return confDir
}

func setupLogging() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
}

func main() {
	app := cli.NewApp()
	app.Name = "dope"

	setupLogging()

	confDir := initConfDir()
	manifest, err := initManifest(confDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app.Commands = []cli.Command{
		{
			Name:    "self-update",
			Aliases: []string{"sup"},
			Usage:   "update dope",
			Action: func(c *cli.Context) error {
				fmt.Println("TODO update dope")
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
					fmt.Println("TODO update", name)
				} else {
					// update all packages
					fmt.Println("TODO update all")
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
						log.Info(repo, "already installed, try update instead")
						return nil
					}

					// install package
					pack, err := installImage(repo)
					if err != nil {
						log.Error(err)
						return err
					}

					if err := manifest.addPack(pack); err != nil {
						log.Error(err)
						return err
					}

					log.Info("Installed", pack.Name, pack.Tag)
				} else {
					err = errors.New("no package name given to install")
					log.Error(err)
					return err
				}
				return nil
			},
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run a package's alias",
			Action: func(c *cli.Context) (err error) {
				notifyIfSelfUpdateAvail()

				if c.NArg() > 0 {
					// run package with args
					name := c.Args()[0]
					args := c.Args()[1:]
					fmt.Println("TODO run", name, args)
				} else {
					err = errors.New("no package name given to run")
				}
				return err
			},
		},
		{
			Name:    "check",
			Aliases: []string{"ch"},
			Usage:   "check for updates to package",
			Action: func(c *cli.Context) (err error) {
				notifyIfSelfUpdateAvail()

				if c.NArg() > 0 {
					// check package for updates
					name := c.Args()[0]
					avail, tag := manifest.checkForUpdate(name)
					if avail {
						fmt.Println("New version", tag, "available for", name)
					} else {
						fmt.Println("No updates available for", name)
					}
				} else {
					err = errors.New("no package name given to check")
				}
				return err
			},
		},
	}

	app.Run(os.Args)
}

func notifyIfSelfUpdateAvail() {
	// TODO better message
	if selfUpdateAvail() {
		fmt.Println("dope update available")
	}
}

func selfUpdateAvail() bool {
	// TODO implement me
	return false
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

	p := &Pack{
		Repo: repo,
		Tag:  tag,
		Name: name,
	}
	return p, nil
}
