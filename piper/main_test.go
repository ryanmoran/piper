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
			"-o", "output-1=/tmp/local-2",
		)
		command.Env = append(os.Environ(), "VAR1=var-1")

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))

		dockerInvocations, err := ioutil.ReadFile(dockerconfig.InvocationsPath)
		Expect(err).NotTo(HaveOccurred())

		dockerCommands := strings.Split(strings.TrimSpace(string(dockerInvocations)), "\n")
		Expect(dockerCommands).To(Equal([]string{
			fmt.Sprintf("%s pull my-image", pathToDocker),
			fmt.Sprintf("%s run --workdir=/tmp/build --env=VAR1=var-1 --volume=/tmp/local-1:/tmp/build/input-1 --volume=/tmp/local-2:/tmp/build/output-1 --tty my-image my-task.sh", pathToDocker),
		}))
	})

	It("runs a concourse task with input image and tag", func() {
		command := exec.Command(pathToPiper,
			"-c", "fixtures/task.yml",
			"-r", "my-image",
			"-t", "my-tag",
			"-i", "input-1=/tmp/local-1",
			"-o", "output-1=/tmp/local-2",
		)
		command.Env = append(os.Environ(), "VAR1=var-1")

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))

		dockerInvocations, err := ioutil.ReadFile(dockerconfig.InvocationsPath)
		Expect(err).NotTo(HaveOccurred())

		dockerCommands := strings.Split(strings.TrimSpace(string(dockerInvocations)), "\n")
		Expect(dockerCommands).To(Equal([]string{
			fmt.Sprintf("%s pull my-image:my-tag", pathToDocker),
			fmt.Sprintf("%s run --workdir=/tmp/build --env=VAR1=var-1 --volume=/tmp/local-1:/tmp/build/input-1 --volume=/tmp/local-2:/tmp/build/output-1 --tty my-image:my-tag my-task.sh", pathToDocker),
		}))
	})

	It("runs a concourse task with complex inputs", func() {
		command := exec.Command(pathToPiper,
			"-c", "fixtures/advanced_task.yml",
			"-i", "input=/tmp/local-1",
			"-o", "output=/tmp/local-2",
			"-p")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))

		dockerInvocations, err := ioutil.ReadFile(dockerconfig.InvocationsPath)
		Expect(err).NotTo(HaveOccurred())

		dockerCommands := strings.Split(strings.TrimSpace(string(dockerInvocations)), "\n")
		Expect(dockerCommands).To(Equal([]string{
			fmt.Sprintf("%s pull my-image:x.y", pathToDocker),
			fmt.Sprintf("%s run --workdir=/tmp/build --privileged --volume=/tmp/local-1:/tmp/build/some/path/input --volume=/tmp/local-2:/tmp/build/some/path/output --tty my-image:x.y my-task.sh", pathToDocker),
		}))
	})

	It("prints the docker commands to stdout, but does not execute them", func() {
		command := exec.Command(pathToPiper,
			"--dry-run",
			"-c", "fixtures/advanced_task.yml",
			"-i", "input=/tmp/local-1",
			"-o", "output=/tmp/local-2",
			"-p")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))

		dockerCommands := strings.Split(strings.TrimSpace(string(session.Out.Contents())), "\n")
		Expect(dockerCommands).To(Equal([]string{
			fmt.Sprintf("%s pull my-image:x.y", pathToDocker),
			fmt.Sprintf("%s run --workdir=/tmp/build --privileged --volume=/tmp/local-1:/tmp/build/some/path/input --volume=/tmp/local-2:/tmp/build/some/path/output --tty my-image:x.y my-task.sh", pathToDocker),
		}))
		_, err = os.Stat(dockerconfig.InvocationsPath)
		Expect(os.IsNotExist(err)).To(BeTrue())
	})

	Context("failure cases", func() {
		Context("when the flag is not passed in", func() {
			It("Print an error and exit with status 1", func() {
				command := exec.Command(pathToPiper)
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				Expect(session.Err.Contents()).To(ContainSubstring("-c is a required flag"))
			})
		})

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
				Expect(session.Err.Contents()).To(ContainSubstring("The following required inputs/outputs are not satisfied: input-1, output-1."))
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
					"-o", "output-1=/tmp/local-2")
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
					"-o", "output-1=/tmp/local-2")
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
					"-o", "output-1=/tmp/local-2")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				Expect(session.Err.Contents()).To(ContainSubstring("failed to run"))
			})
		})
	})
})
