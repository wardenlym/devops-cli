package main

import (
	"embed"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/wardenlym/devops-cli/cmd"
)

//go:embed infra
var infra embed.FS

// VERSION gets overridden at build time using -X main.VERSION=$VERSION
var VERSION = "dev"
var released = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+$`)

func main() {
	cmd.Infra = infra
	logrus.SetOutput(colorable.NewColorableStdout())

	if err := mainErr(); err != nil {
		logrus.Fatal(err)
	}
}

func mainErr() error {
	app := cli.NewApp()
	app.Name = "devops-cli"
	app.Version = VERSION
	app.Usage = "devops-cli"
	app.Before = func(ctx *cli.Context) error {
		if ctx.GlobalBool("quiet") {
			logrus.SetOutput(ioutil.Discard)
		} else {
			if ctx.GlobalBool("debug") {
				logrus.SetLevel(logrus.DebugLevel)
				logrus.Debugf("Loglevel set to [%v]", logrus.DebugLevel)
			}
			if ctx.GlobalBool("trace") {
				logrus.SetLevel(logrus.TraceLevel)
				logrus.Tracef("Loglevel set to [%v]", logrus.TraceLevel)
			}
		}
		// if released.MatchString(app.Version) {
		// 	metadata.RKEVersion = app.Version
		// 	return nil
		// }
		// logrus.Warnf("This is not an officially supported version (%s) of RKE. Please download the latest official release at https://github.com/rancher/rke/releases", app.Version)
		return nil
	}
	app.Author = "wardenlym"
	app.Email = ""
	app.Commands = []cli.Command{
		cmd.InitCommand(),
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug,d",
			Usage: "Debug logging",
		},
		cli.BoolFlag{
			Name:  "quiet,q",
			Usage: "Quiet mode, disables logging and only critical output will be printed",
		},
		cli.BoolFlag{
			Name:  "trace",
			Usage: "Trace logging",
		},
	}
	return app.Run(os.Args)
}
