package piper

import (
	"fmt"
	"path/filepath"
	"strings"
)

const VolumeMountPoint = "/tmp/build"

type VolumeMountBuilder struct{}

func (b VolumeMountBuilder) Build(inputs []VolumeMount, pairs []string) ([]DockerVolumeMount, error) {
	pairsMap := make(map[string]string)

	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("could not parse input %q. must be of form <input-name>=<input-location>", pair)
		}

		pairsMap[parts[0]] = parts[1]
	}

	var mounts []DockerVolumeMount
	for _, input := range inputs {
		inputLocation, ok := pairsMap[input.Name]
		if !ok {
			return nil, fmt.Errorf("input %q is not satisfied. please include an input in command arguments", input.Name)
		}
		var mountPoint string
		if input.Path == "" {
			mountPoint = filepath.Join(VolumeMountPoint, input.Name)
		} else {
			mountPoint = filepath.Join(VolumeMountPoint, input.Path)
		}

		mounts = append(mounts, DockerVolumeMount{
			LocalPath:  inputLocation,
			RemotePath: filepath.Clean(mountPoint),
		})
	}

	return mounts, nil
}
