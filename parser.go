package piper

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type VolumeMount struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type Run struct {
	Path string   `yaml:"path"`
	args []string `yaml:"args"`
	dir  string   `yaml:"dir"`
}

type Task struct {
	Image  string `yaml:"image"`
	Run    Run
	Inputs []VolumeMount
	Params map[string]string
}

type Parser struct{}

func (p Parser) Parse(path string) (Task, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return Task{}, err
	}

	var task Task

	err = yaml.Unmarshal(contents, &task)
	if err != nil {
		return Task{}, err
	}

	task.Image = strings.TrimPrefix(task.Image, "docker:///")

	return task, nil
}
