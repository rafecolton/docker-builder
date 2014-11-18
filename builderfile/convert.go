package builderfile

import (
	"time"
)

/*
Convert0to1 converts the deprecated unitconfig version 0 to version 1.  It
also prints out a deprecation warning message.
*/
func Convert0to1(zero *UnitConfig) (*UnitConfig, error) {
	logger.Warn(versionZeroWarningMessage)
	time.Sleep(sleepTime * time.Second)

	ret := &UnitConfig{
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
