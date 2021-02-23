package cmd

import (
	"io/fs"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	comments = `# this is some comments,
# hold this as todo.`
)

func InitCommand() cli.Command {
	return cli.Command{
		Name:   "init",
		Usage:  "Init .devops-cli working directory",
		Action: initWorkingDirectory,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "path,p",
				Usage: "Path of the working directory",
				Value: func() string {
					user, err := user.Current()
					if err != nil {
						log.Fatalf(err.Error())
					}
					return filepath.Join(user.HomeDir, WorkingDirName)
				}(),
			},
		},
	}
}

func initWorkingDirectory(ctx *cli.Context) error {
	err := os.Mkdir(ctx.String("path"), 0755)
	if err != nil {
		// if os.IsExist(err) {
		// 	err = os.Remove(ctx.String("path"))
		// 	if err != nil {
		// 		return err
		// 	}
		// 	err = os.Mkdir(ctx.String("path"), 0777)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		return err
	}

	err = fs.WalkDir(Infra, ".", func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if d.IsDir() {
			e = os.MkdirAll(filepath.Join(ctx.String("path"), path), 0755)
		} else {
			b, e2 := fs.ReadFile(Infra, path)
			if e2 != nil {
				logrus.Error(e2)
			}

			e = os.WriteFile(filepath.Join(ctx.String("path"), path), b, 0644)
		}
		if e != nil {
			logrus.Error(e)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
