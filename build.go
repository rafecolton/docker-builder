package main

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/rafecolton/docker-builder/conf"
	"github.com/winchman/builder-core"
	"github.com/winchman/builder-core/unit-config"

	"github.com/codegangsta/cli"
	"github.com/onsi/gocleanup"
)

func build(c *cli.Context) {
	builderfile := c.Args().First()
	if builderfile == "" {
		builderfile = "Bobfile"
	}

	unitConfig, err := unitconfig.ReadFromFile("./"+builderfile, unitconfig.TOML, unitconfig.YAML)
	if err != nil {
		if c.Bool("force") {
			if err := forceBuild(); err != nil {
				Logger.Warn(err.Error())
			}
		}
		gocleanup.Exit(0)
	}

	globals := unitconfig.ConfigGlobals{
		SkipPush: c.Bool("skip-push") || conf.Config.SkipPush,
		CfgUn:    conf.Config.CfgUn,
		CfgPass:  conf.Config.CfgPass,
		CfgEmail: conf.Config.CfgEmail,
	}

	unitConfig.SetGlobals(globals)

	if err := runner.RunBuildSynchronously(runner.Options{
		UnitConfig: unitConfig,
		ContextDir: os.Getenv("PWD"),
		LogLevel:   Logger.Level,
	}); err != nil {
		exitErr(1, "unable to build", err)
	}

}

func forceBuild() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	basename := path.Base(pwd)
	args := []string{"docker", "build", "-t", basename, "."}

	cmd := exec.Command("docker")
	cmd.Args = args
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	Logger.Info("running command --> " + strings.Join(args, " "))
	return cmd.Run()
}
