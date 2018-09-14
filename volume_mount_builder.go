package piper

import (
	"fmt"
	"path/filepath"
	"strings"
)

const VolumeMountPoint = "/tmp/build"

type VolumeMountBuilder struct{}

func (b VolumeMountBuilder) Build(resources []VolumeMount, inputs, outputs []string) ([]DockerVolumeMount, error) {
	pairsMap := make(map[string]string)

	for _, input := range inputs {
		parts := strings.Split(input, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("could not parse input %q. must be of form <input-name>=<input-location>", input)
		}

		pairsMap[parts[0]] = parts[1]
	}

	for _, output := range outputs {
		parts := strings.Split(output, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("could not parse output %q. must be of form <output-name>=<output-location>", output)
		}

		pairsMap[parts[0]] = parts[1]
	}

	var mounts []DockerVolumeMount
	var missingResources []string
	for _, resource := range resources {
		resourceLocation, ok := pairsMap[resource.Name]
		if !ok {
			if !resource.Optional {
				missingResources = append(missingResources, resource.Name)
			}
			continue
		}
		var mountPoint string
		if resource.Path == "" {
			mountPoint = filepath.Join(VolumeMountPoint, resource.Name)
		} else {
			mountPoint = filepath.Join(VolumeMountPoint, resource.Path)
		}

		mounts = append(mounts, DockerVolumeMount{
			LocalPath:  resourceLocation,
			RemotePath: filepath.Clean(mountPoint),
		})
	}
	if len(missingResources) != 0 {
		return nil, fmt.Errorf("The following required inputs/outputs are not satisfied: %s.", strings.Join(missingResources, ", "))
	}

	return mounts, nil
}
