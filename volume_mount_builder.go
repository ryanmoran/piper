package piper

import (
	"fmt"
	"strings"
)

const VolumeMountPoint = "/tmp/build"

type VolumeMountBuilder struct{}

func (b VolumeMountBuilder) Build(inputs, pairs []string, folderType string) ([]DockerVolumeMount, error) {
	pairsMap := make(map[string]string)

	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("could not parse %s %q. must be of form <%s-name>=<%s-location>",folderType, pair, folderType, folderType)
		}

		pairsMap[parts[0]] = parts[1]
	}

	var mounts []DockerVolumeMount
	for _, input := range inputs {
		inputLocation, ok := pairsMap[input]
		if !ok {
			return nil, fmt.Errorf("%s %q is not satisfied. please include an %s in command arguments",folderType, input, folderType)
		}

		mounts = append(mounts, DockerVolumeMount{
			LocalPath:  inputLocation,
			RemotePath: fmt.Sprintf("%s/%s", VolumeMountPoint, input),
		})
	}

	return mounts, nil
}
