package main_test

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestPiperExecutable(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "piper/piper")
}

var (
	pathToPiper  string
	pathToDocker string
)

var _ = BeforeSuite(func() {
	var err error
	pathToPiper, err = gexec.Build("github.com/ryanmoran/piper/piper")
	Expect(err).NotTo(HaveOccurred())

	pathToDocker, err = gexec.Build("github.com/ryanmoran/piper/fakes/docker")
	Expect(err).NotTo(HaveOccurred())

	os.Setenv("PATH", fmt.Sprintf("%s:%s", filepath.Dir(pathToDocker), os.Getenv("PATH")))
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
