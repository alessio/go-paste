package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-ini/ini"
	"github.com/urfave/cli/v2"

	"github.com/bearbin/go-paste/fpaste"
	"github.com/bearbin/go-paste/pastebin"
	"github.com/bearbin/go-paste/ubuntu"
)

var errUnknownService = errors.New("unknown paste service")
var iniFile *ini.File
var iniFilename = "paste.ini"

func loadConfig() {
	var err error

	iniFile, err = ini.Load(iniFilename)
	if err != nil {
		iniFile = ini.Empty()
	}
}

func main() {
	app := cli.NewApp()
	app.Usage = "get and put pastes from pastebin and other paste sites."
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "service, s", Value: "ubuntu", Usage: "the pastebin service to use"},
	}
	app.Commands = []*cli.Command{
		{
			Name:    "put",
			Usage:   "put a paste",
			Aliases: []string{"p"},
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "id", Usage: "return the paste id not the url"},
				&cli.StringFlag{Name: "title, t", Value: "", Usage: "the title for the paste"},
			},
			Action: func(c *cli.Context) error {
				srv, err := convertService(c.String("service"))
				if err != nil {
					return err
				}

				var text []byte
				if c.Args().First() == "-" || c.Args().First() == "" {
					text, err = ioutil.ReadAll(os.Stdin)
				} else {
					text, err = ioutil.ReadFile(c.Args().First())
				}

				if err != nil {
					return err
				}

				code, err := srv.Put(string(text), c.String("title"))
				if err != nil {
					return err
				}

				if !c.Bool("id") {
					code = srv.WrapID(code)
				}

				_, _ = fmt.Fprintln(app.Writer, code)

				return nil
			},
		},
		{
			Name:    "get",
			Usage:   "get a paste from its url",
			Aliases: []string{"g"},
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "id", Usage: "get a paste from its ID instead of its URL"},
			},
			Action: func(c *cli.Context) error {
				srv, err := convertService(c.String("service"))
				if err != nil {
					fmt.Fprintf(app.ErrWriter, "ERROR: %v\n", err.Error())
					os.Exit(1)
				}

				var id string
				if c.Bool("id") {
					id = c.Args().First()
				} else {
					id = srv.StripURL(c.Args().First())
				}

				text, err := srv.Get(id)
				if err != nil {
					return err
				}

				_, _ = fmt.Fprintln(app.Writer, text)

				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(app.ErrWriter, "%v\n", err)

		os.Exit(1)
	}
}

func convertService(srv string) (service, error) {
	loadConfig()

	switch {
	case srv == "pastebin" || srv == "pastebin.com" || srv == "http://pastebin.com":
		pastebinConfig := pastebin.DefaultConfig()
		cfgFileSection := iniFile.Section(pastebin.Name)

		if err := cfgFileSection.MapTo(pastebinConfig); err != nil {
			panic(err)
		}

		return pastebin.New(pastebinConfig), nil
	case srv == "fpaste" || srv == "fpaste.org" || srv == "http://fpaste.org":
		return fpaste.Fpaste{}, nil
	case srv == "ubuntu" || srv == "paste.ubuntu.com":
		pastebinConfig := ubuntu.DefaultConfig()
		cfgFileSection := iniFile.Section(pastebin.Name)

		if err := cfgFileSection.MapTo(pastebinConfig); err != nil {
			panic(err)
		}

		return ubuntu.New(pastebinConfig), nil
	}

	return nil, errUnknownService
}
