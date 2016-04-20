package piper

import (
	"io/ioutil"
	"strings"

	"github.com/cloudfoundry/cf-release/src/dea-hm-workspace/src/dea_next/go/src/github.com/go-yaml/yaml"
)

type Config struct {
	Image   string
	Command string
	Inputs  []string
	Params  map[string]string
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
		Params map[string]string `yaml:"params"`
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
		Params:  task.Params,
	}, nil
}
