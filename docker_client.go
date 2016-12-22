package piper

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
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

func (c DockerClient) Pull(image string, dryRun bool) error {
	args := append(c.Command.Args, "pull", image)

	if dryRun {
		fmt.Fprintln(c.Stdout, strings.Join(args, " "))
		return nil
	}

	c.Command.Args = args
	c.Command.Stdout = c.Stdout
	c.Command.Stderr = c.Stderr

	err := c.Command.Run()
	if err != nil {
		return err
	}

	return nil
}

func (c DockerClient) Run(command, image string, envVars []DockerEnv, mounts []DockerVolumeMount, privileged bool, dryRun bool) error {
	c.Command.Args = append(c.Command.Args, "run", fmt.Sprintf("--workdir=%s", VolumeMountPoint))

	if privileged {
		c.Command.Args = append(c.Command.Args, "--privileged")
	}

	for _, envVar := range envVars {
		c.Command.Args = append(c.Command.Args, envVar.String())
	}

	for _, mount := range mounts {
		c.Command.Args = append(c.Command.Args, mount.String())
	}

	c.Command.Args = append(c.Command.Args, image, command)

	if dryRun {
		fmt.Fprintln(c.Stdout, strings.Join(c.Command.Args, " "))
		return nil
	}

	c.Command.Stdout = c.Stdout
	c.Command.Stderr = c.Stderr

	err := c.Command.Run()
	if err != nil {
		return err
	}

	return nil
}
