package piper

import (
	"fmt"
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

type ImageResourceSource struct {
	Repository string
	Tag        string
}

func (i ImageResourceSource) String() string {
	if "" != i.Tag {
		return fmt.Sprintf("%s:%s", i.Repository, i.Tag)
	}
	return i.Repository
}

type ImageResource struct {
	Source ImageResourceSource
}

type Task struct {
	Image         string `yaml:"image"`
	Run           Run
	Inputs        []VolumeMount
	Outputs       []VolumeMount
	Params        map[string]string
	ImageResource ImageResource `yaml:"image_resource"`
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
	if task.ImageResource.Source.Repository != "" {
		task.Image = task.ImageResource.Source.String()
	} else {
		task.Image = strings.TrimPrefix(task.Image, "docker:///")
	}

	return task, nil
}
