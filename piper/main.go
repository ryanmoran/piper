package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/ryanmoran/piper"
)

func main() {
	var (
		taskFilePath string
		inputPairs   ResourcePairs
		outputPairs  ResourcePairs
		privileged   bool
		dryRun       bool
	)

	flag.StringVar(&taskFilePath, "c", "", "path to the task configuration file")
	flag.Var(&inputPairs, "i", "<input-name>=<input-location>")
	flag.Var(&outputPairs, "o", "<output-name>=<output-location>")
	flag.BoolVar(&privileged, "p", false, "run the task with full privileges")
	flag.BoolVar(&dryRun, "dry-run", false, "prints the docker commands without running them")

	flag.Parse()

	var errors []string
	if len(taskFilePath) == 0 {
		errors = append(errors, fmt.Sprintf(" -c is a required flag"))
	}

	if len(errors) > 0 {
		fmt.Fprintln(os.Stderr, "Errors:")
		for _, err := range errors {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Fprintln(os.Stderr, "\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	taskConfig, err := piper.Parser{}.Parse(taskFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	var resources []piper.VolumeMount
	resources = append(resources, taskConfig.Inputs...)
	resources = append(resources, taskConfig.Outputs...)

	volumeMounts, err := piper.VolumeMountBuilder{}.Build(resources, inputPairs, outputPairs)
	if err != nil {
		log.Fatalln(err)
	}

	envVars := piper.EnvVarBuilder{}.Build(os.Environ(), taskConfig.Params)

	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		log.Fatalln(err)
	}

	dockerClient := piper.DockerClient{
		Command: exec.Command(dockerPath),
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}
	err = dockerClient.Pull(taskConfig.Image, dryRun)
	if err != nil {
		log.Fatalln(err)
	}

	dockerClient = piper.DockerClient{
		Command: exec.Command(dockerPath),
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}

	err = dockerClient.Run(taskConfig.Run.Path, taskConfig.Image, envVars, volumeMounts, privileged, dryRun)
	if err != nil {
		log.Fatalln(err)
	}
}

type ResourcePairs []string

func (p *ResourcePairs) Set(resource string) error {
	*p = append(*p, resource)
	return nil
}

func (p *ResourcePairs) String() string {
	return fmt.Sprint(*p)
}
