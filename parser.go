package piper

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Image   string
	Command string
	Inputs  []string
}

type Parser struct{}

func (p Parser) Parse(path string) (Config, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var task struct {
		Image string `yaml:"image"`
		Run   struct {
			Path string `yaml:"path"`
		} `yaml:"run"`
		Inputs []struct {
			Name string `yaml:"name"`
		} `yaml:"inputs"`
	}

	err = yaml.Unmarshal(contents, &task)
	if err != nil {
		return Config{}, err
	}

	var inputs []string
	for _, input := range task.Inputs {
		inputs = append(inputs, input.Name)
	}

	return Config{
		Image:   strings.TrimPrefix(task.Image, "docker:///"),
		Command: task.Run.Path,
		Inputs:  inputs,
	}, nil
}
