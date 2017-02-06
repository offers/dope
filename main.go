package main

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "dope"

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
					name := c.Args()[0]
					fmt.Println("TODO install", name)
				} else {
					err = errors.New("no package name given to install")
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
					fmt.Println("TODO check", name)
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
