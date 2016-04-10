package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ryanmoran/piper/fakes/docker/dockerconfig"
)

var failPull, failRun bool

func main() {
	command := strings.Join(os.Args, " ")

	if failPull && strings.Contains(command, "docker pull") {
		log.Fatalln("failed to pull")
	}

	if failRun && strings.Contains(command, "docker run") {
		log.Fatalln("failed to run")
	}

	err := os.MkdirAll(filepath.Dir(dockerconfig.InvocationsPath), 0755)
	if err != nil {
		log.Fatalln(err)
	}

	invocations, err := os.OpenFile(dockerconfig.InvocationsPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = invocations.WriteString(command + "\n")
	if err != nil {
		log.Fatalln(err)
	}

	err = invocations.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
