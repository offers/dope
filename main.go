package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/urfave/cli.v1"
)

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

func main() {
	app := cli.NewApp()
	app.Name = "dope"

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
					// install package
					image := c.Args()[0]
					p := newPack(image)
					if err := manifest.addPack(p); err != nil {
						fmt.Println(err)
						return err
					}
					fmt.Println("installed", p.Name)
				} else {
					fmt.Println(err)
					err = errors.New("no package name given to install")
				}
				return err
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
