package builderfile

import (
	"time"

	"github.com/wsxiaoys/terminal/color"
)

const sleepTime = 30

var versionZeroWarninMessage = color.Sprintf(`@{y!}WARNING, PLEASE READ:@{|y}

Uh oh, it looks like you are using a Bobfile with @{!y}a deprecated format@{|y}.

You may want to @{!y}update@{|y} before proceeding.

You should consider @{!y}pressing ^C@{|y} and @{!y}updating your Bobfile@{|y} before proceeding.

For more information, see @{!y}https://github.com/modcloth/docker-builder@{|y}

Sleeping for 30 seconds before building...
@{|}`)

/*
Convert0to1 converts the deprecated builderfile version 0 to version 1.  It
also prints out a deprecation warning message.
*/
func Convert0to1(zero *Builderfile) (*Builderfile, error) {
	logger.Warn(versionZeroWarninMessage)
	time.Sleep(sleepTime * time.Second)

	ret := &Builderfile{
		Version:          1,
		Docker:           zero.Docker,
		ContainerGlobals: &ContainerSection{},
		ContainerArr:     []*ContainerSection{},
	}

	if zero.Containers == nil || len(zero.Containers) == 0 {
		return ret, nil
	}

	if globals, ok := zero.Containers["global"]; ok {
		*ret.ContainerGlobals = globals
	}

	for key, container := range zero.Containers {
		if key == "global" {
			continue
		}

		ret.ContainerArr = append(ret.ContainerArr, &container)
	}

	return ret, nil
}
