package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/modcloth/docker-builder/analyzer"

	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func initialize(c *cli.Context) {
	dir := c.Args().First()
	if dir == "" {
		dir = "."
	}

	file, err := analyzer.ParseAnalysisFromDir(dir)
	if err != nil {
		exitErr(1, "unable to create Bobfile", err)
	}

	bobfilePath := filepath.Join(dir, "Bobfile")

	//no error when stating, file already exists, rename with timestamp
	if _, err := os.Stat(bobfilePath); err == nil {
		bobfilePath = fmt.Sprintf("%s.%d", bobfilePath, int32(time.Now().Unix()))
	}

	outfile, err := os.Create(bobfilePath)
	if err != nil {
		exitErr(86, "unable to create output file", map[string]interface{}{"output_file": bobfilePath, "error": err})
	}
	defer outfile.Close()

	encoder := toml.NewEncoder(outfile)
	if err = encoder.Encode(file); err != nil {
		exitErr(123, "unable to write to output file", map[string]interface{}{"output_file": bobfilePath, "error": err})
	}

	vimFtComment := []byte("\n\n# vim:ft=toml")
	if _, err := outfile.Write(vimFtComment); err != nil {
		exitErr(127, "unable to write to output file", map[string]interface{}{"output_file": bobfilePath, "error": err})
	}

	Logger.WithFields(logrus.Fields{"output_file": bobfilePath}).Info("successfully initialized")
}
