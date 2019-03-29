package piper

import (
	"fmt"
	"os/user"
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

		expandedPath, err := expandUser(parts[1])
		if err != nil {
			return nil, err
		}
		pairsMap[parts[0]] = expandedPath
	}

	for _, output := range outputs {
		parts := strings.Split(output, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("could not parse output %q. must be of form <output-name>=<output-location>", output)
		}

		expandedPath, err := expandUser(parts[1])
		if err != nil {
			return nil, err
		}
		pairsMap[parts[0]] = expandedPath
	}

	var mounts []DockerVolumeMount
	var missingResources []string
	for _, resource := range resources {
		if resource.Name == "" && resource.Path != "" {
			mountPoint := filepath.Join(VolumeMountPoint, resource.Path)

			mounts = append(mounts, DockerVolumeMount{
				LocalPath:  "/tmp",
				RemotePath: filepath.Clean(mountPoint),
			})
			continue
		}

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

func expandUser(path string) (string, error) {
	if !strings.HasPrefix(path, "~/") {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	dir := usr.HomeDir
	return filepath.Join(dir, path[2:]), nil
}
