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
		inputPairs   StringPairs
		outputPairs   StringPairs
	)

	flag.StringVar(&taskFilePath, "c", "", "path to the task configuration file")
	flag.Var(&inputPairs, "i", "<input-name>=<input-location>")
	flag.Var(&outputPairs, "o", "<output-name>=<output-location>")

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

	volumeMounts, err := piper.VolumeMountBuilder{}.Build(taskConfig.Inputs, inputPairs, "input")
	if err != nil {
		log.Fatalln(err)
	}

  fmt.Println(taskConfig.Outputs)
	outputVolumeMounts, err := piper.VolumeMountBuilder{}.Build(taskConfig.Outputs, outputPairs, "output")
	if err != nil {
		log.Fatalln(err)
	}
	volumeMounts = append(volumeMounts, outputVolumeMounts...)

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
	err = dockerClient.Pull(taskConfig.Image)
	if err != nil {
		log.Fatalln(err)
	}

	dockerClient = piper.DockerClient{
		Command: exec.Command(dockerPath),
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}

	err = dockerClient.Run(taskConfig.Command, taskConfig.Image, envVars, volumeMounts)
	if err != nil {
		log.Fatalln(err)
	}
}

type StringPairs []string

func (i *StringPairs) Set(input string) error {
	*i = append(*i, input)
	return nil
}

func (i *StringPairs) String() string {
	return fmt.Sprint(*i)
}
