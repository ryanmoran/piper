package piper

import (
	"fmt"
	"io"
	"os/exec"
)

type DockerVolumeMount struct {
	LocalPath  string
	RemotePath string
}

func (m DockerVolumeMount) String() string {
	return fmt.Sprintf("--volume=%s:%s", m.LocalPath, m.RemotePath)
}

type DockerEnv struct {
	Key   string
	Value string
}

func (e DockerEnv) String() string {
	return fmt.Sprintf("--env=%s=%s", e.Key, e.Value)
}

type DockerClient struct {
	Command *exec.Cmd
	Stdout  io.Writer
	Stderr  io.Writer
}

func (c DockerClient) Pull(image string) error {
	c.Command.Args = append(c.Command.Args, "pull", image)
	c.Command.Stdout = c.Stdout
	c.Command.Stderr = c.Stderr

	err := c.Command.Run()
	if err != nil {
		return err
	}

	return nil
}

func (c DockerClient) Run(command, image string, envVars []DockerEnv, mounts []DockerVolumeMount) error {
	args := []string{
		"run",
		fmt.Sprintf("--workdir=%s", VolumeMountPoint),
	}

	for _, envVar := range envVars {
		args = append(args, envVar.String())
	}

	for _, mount := range mounts {
		args = append(args, mount.String())
	}

	args = append(args, image, command)

	c.Command.Args = append(c.Command.Args, args...)
	c.Command.Stdout = c.Stdout
	c.Command.Stderr = c.Stderr

	err := c.Command.Run()
	if err != nil {
		return err
	}

	return nil
}
