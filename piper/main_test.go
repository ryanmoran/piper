package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/onsi/gomega/gexec"
	"github.com/ryanmoran/piper/fakes/docker/dockerconfig"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Piper", func() {
	BeforeEach(func() {
		err := os.RemoveAll(dockerconfig.InvocationsPath)
		Expect(err).NotTo(HaveOccurred())
	})

	It("runs a concourse task", func() {
		command := exec.Command(pathToPiper,
			"-c", "fixtures/task.yml",
			"-i", "input-1=/tmp/local-1",
			"-i", "input-2=/tmp/local-2")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))

		dockerInvocations, err := ioutil.ReadFile(dockerconfig.InvocationsPath)
		Expect(err).NotTo(HaveOccurred())

		dockerCommands := strings.Split(strings.TrimSpace(string(dockerInvocations)), "\n")
		Expect(dockerCommands).To(Equal([]string{
			fmt.Sprintf("%s pull my-image", pathToDocker),
			fmt.Sprintf("%s run --workdir=/tmp/build --volume=/tmp/local-1:/tmp/build/input-1 --volume=/tmp/local-2:/tmp/build/input-2 my-image my-task.sh", pathToDocker),
		}))
	})

	Context("failure cases", func() {
		Context("when the task file does not exist", func() {
			It("prints an error and exits 1", func() {
				command := exec.Command(pathToPiper, "-c", "no-such-file")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				Expect(session.Err.Contents()).To(ContainSubstring("no such file or directory"))
			})
		})

		Context("when inputs are missing", func() {
			It("prints an error and exits 1", func() {
				command := exec.Command(pathToPiper, "-c", "fixtures/task.yml")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				Expect(session.Err.Contents()).To(ContainSubstring("input \"input-1\" is not satisfied"))
			})
		})

		Context("when docker cannot be found on the $PATH", func() {
			var path string

			BeforeEach(func() {
				path = os.Getenv("PATH")
				os.Setenv("PATH", "")
			})

			AfterEach(func() {
				os.Setenv("PATH", path)
			})

			It("prints an error and exits 1", func() {
				command := exec.Command(pathToPiper,
					"-c", "fixtures/task.yml",
					"-i", "input-1=/tmp/local-1",
					"-i", "input-2=/tmp/local-2")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				Expect(session.Err.Contents()).To(ContainSubstring("executable file not found in $PATH"))
			})
		})

		Context("when docker fails to pull the image", func() {
			var pathToBadDocker, path string

			BeforeEach(func() {
				var err error
				pathToBadDocker, err = gexec.Build("github.com/ryanmoran/piper/fakes/docker", "-tags", "fail_pull")
				Expect(err).NotTo(HaveOccurred())

				path = os.Getenv("PATH")
				os.Setenv("PATH", fmt.Sprintf("%s:%s", filepath.Dir(pathToBadDocker), os.Getenv("PATH")))
			})

			AfterEach(func() {
				os.Setenv("PATH", path)
			})

			It("prints an error and exits 1", func() {
				command := exec.Command(pathToPiper,
					"-c", "fixtures/task.yml",
					"-i", "input-1=/tmp/local-1",
					"-i", "input-2=/tmp/local-2")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				Expect(session.Err.Contents()).To(ContainSubstring("failed to pull"))
			})
		})

		Context("when docker fails to run the command", func() {
			var pathToBadDocker, path string

			BeforeEach(func() {
				var err error
				pathToBadDocker, err = gexec.Build("github.com/ryanmoran/piper/fakes/docker", "-tags", "fail_run")
				Expect(err).NotTo(HaveOccurred())

				path = os.Getenv("PATH")
				os.Setenv("PATH", fmt.Sprintf("%s:%s", filepath.Dir(pathToBadDocker), os.Getenv("PATH")))
			})

			AfterEach(func() {
				os.Setenv("PATH", path)
			})

			It("prints an error and exits 1", func() {
				command := exec.Command(pathToPiper,
					"-c", "fixtures/task.yml",
					"-i", "input-1=/tmp/local-1",
					"-i", "input-2=/tmp/local-2")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				Expect(session.Err.Contents()).To(ContainSubstring("failed to run"))
			})
		})
	})
})
