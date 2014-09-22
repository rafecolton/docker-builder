package main

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/rafecolton/docker-builder/builder"

	"github.com/codegangsta/cli"
)

func build(c *cli.Context) {
	builder.SkipPush = c.Bool("skip-push")
	builderfile := c.Args().First()
	if builderfile == "" {
		builderfile = "Bobfile"
	}

	bob, err := builder.NewBuilder(Logger, true)
	if err != nil {
		exitErr(61, "unable to build", err)
	}

	config, err := builder.NewTrustedFilePath(builderfile, ".")
	if err != nil {
		exitErr(1, "unable to create build config", err)
	}

	if err := bob.Build(config); err != nil {
		if builder.IsSanitizeError(err) {
			if c.Bool("force") {
				pwd, err := os.Getwd()
				if err != nil {
					exitErr(1, "unable to get cwd", err)
				}

				basename := path.Base(pwd)
				args := []string{"docker", "build", "-t", basename, "."}

				cmd := exec.Command("docker")
				cmd.Args = args
				cmd.Stdout = os.Stdout
				cmd.Stdin = os.Stdin

				Logger.Info("running command --> " + strings.Join(args, " "))
				if err = cmd.Run(); err != nil {
					exitErr(1, "docker build failed", err)
				}
				return
			}

			exitErr(err.ExitCode(), "unable to build", map[string]interface{}{
				"error":    err,
				"filename": err.(*builder.SanitizeError).Filename,
			})
		} else {
			exitErr(err.ExitCode(), "unable to build", err)
		}
	}
}
