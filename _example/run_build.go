package main

import (
	"fmt"
	"os"

	"github.com/sylphon/build-runner"
	"github.com/sylphon/build-runner/builderfile"
)

var example = &builderfile.UnitConfig{
	Version: 1,
	ContainerArr: []*builderfile.ContainerSection{
		&builderfile.ContainerSection{
			Name:       "app",
			Dockerfile: "Dockerfile",
			Registry:   "quay.io/rafecolton",
			Project:    "build-runner-test",
			Tags:       []string{"latest", "git:sha", "git:tag", "git:branch"},
			SkipPush:   true,
		},
	},
}

func main() {
	if err := buildrunner.RunBuild(example, os.Getenv("GOPATH")+"/src/github.com/rafecolton/docker-builder"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
