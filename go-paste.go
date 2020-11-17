package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-ini/ini"
	"github.com/urfave/cli/v2"

	"github.com/bearbin/go-paste/debian"
	"github.com/bearbin/go-paste/pastebin"
	"github.com/bearbin/go-paste/ubuntu"
)

const (
	iniFilename = ".paste.ini"
)

var (
	errUnknownService = errors.New("unknown paste service")
	iniFile           *ini.File
)

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

func loadConfig() {
	var err error

	userHome, _ := os.UserHomeDir() // discard error, file won't be loaded if homedir is not defined
	iniFile, err = ini.LooseLoad(filepath.Join(userHome, iniFilename), iniFilename)

	if err != nil {
		iniFile = ini.Empty()
	}
}

type config struct {
	Pastebin *pastebin.Config `ini:"pastebin.com"`
	Ubuntu   *ubuntu.Config   `ini:"paste.ubuntu.com"`
	Debian   *debian.Config   `ini:"paste.debian.net"`
}

func defaultConfig() *config {
	return &config{
		Pastebin: pastebin.DefaultConfig(),
		Ubuntu:   ubuntu.DefaultConfig(),
		Debian:   debian.DefaultConfig(),
	}
}

func convertService(srv string) (service, error) {
	loadConfig()

	cfg := defaultConfig()
	if err := iniFile.MapTo(cfg); err != nil {
		panic(err)
	}

	switch {
	case findString(srv, "pastebin", "pastebin.com", "http://pastebin.com", "https://pastebin.com") != -1:
		return pastebin.New(cfg.Pastebin), nil
	case findString(srv, "ubuntu", "paste.ubuntu.com", "http://paste.ubuntu.com", "https://paste.ubuntu.com") != -1:
		return ubuntu.New(cfg.Ubuntu), nil
	case findString(srv, "debian", "paste.debian.net", "http://paste.debian.net", "https://paste.debian.net") != -1:
		return debian.New(cfg.Debian), nil
	}

	return nil, errUnknownService
}

func findString(s string, arg ...string) int {
	for index, cur := range arg {
		if cur == s {
			return index
		}
	}

	return -1
}
