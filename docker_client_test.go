package piper_test

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ryanmoran/piper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DockerClient", func() {
	var (
		client piper.DockerClient
		stdout *bytes.Buffer
	)

	BeforeEach(func() {
		stdout = bytes.NewBuffer([]byte{})
		command := exec.Command("echo")

		client = piper.DockerClient{
			Command: command,
			Stdout:  stdout,
		}
	})

	Describe("Pull", func() {
		It("pulls the specified docker image", func() {
			err := client.Pull("some-image", false)
			Expect(err).NotTo(HaveOccurred())

			Expect(stdout.String()).To(Equal("pull some-image\n"))
		})

		It("prints the docker command without running it", func() {
			err := client.Pull("some-image", true)
			Expect(err).NotTo(HaveOccurred())

			Expect(stdout.String()).To(Equal("echo pull some-image\n"))
		})

		Context("failure cases", func() {
			Context("when the executable cannot be found", func() {
				It("returns an error", func() {
					client = piper.DockerClient{
						Command: exec.Command("no-such-executable"),
						Stdout:  stdout,
					}
					err := client.Pull("some-image", false)
					Expect(err).To(MatchError(ContainSubstring("executable file not found in $PATH")))
				})
			})
		})
	})

	Describe("Run", func() {
		It("runs the command with the given volume mounts, and environment", func() {
			err := client.Run([]string{"my-task.sh", "-my-arg1", "-my-arg2"}, "my-image", []piper.DockerEnv{
				{
					Key:   "VAR1",
					Value: "var-1",
				},
				{
					Key:   "VAR2",
					Value: "var-2",
				},
			}, []piper.DockerVolumeMount{
				{
					LocalPath:  "/some/local/path-1",
					RemotePath: "/some/remote/path-1",
				},
				{
					LocalPath:  "/some/local/path-2",
					RemotePath: "/some/remote/path-2",
				},
			}, false, false, false)
			Expect(err).NotTo(HaveOccurred())

			args := []string{
				"run",
				"--workdir=/tmp/build",
				"--env=VAR1=var-1",
				"--env=VAR2=var-2",
				"--volume=/some/local/path-1:/some/remote/path-1",
				"--volume=/some/local/path-2:/some/remote/path-2",
				"-tty",
				"my-image",
				"my-task.sh",
				"-my-arg1",
				"-my-arg2",
			}

			Expect(stdout.String()).To(Equal(strings.Join(args, " ") + "\n"))
		})

		It("runs the command in privileged mode", func() {
			err := client.Run([]string{"my-task.sh"}, "my-image",
				[]piper.DockerEnv{},
				[]piper.DockerVolumeMount{}, true, false, false)
			Expect(err).NotTo(HaveOccurred())

			args := []string{
				"run",
				"--workdir=/tmp/build",
				"--privileged",
				"-tty",
				"my-image",
				"my-task.sh",
			}

			Expect(stdout.String()).To(Equal(strings.Join(args, " ") + "\n"))
		})

		It("runs the command with --rm argument", func() {
			err := client.Run([]string{"my-task.sh"}, "my-image",
				[]piper.DockerEnv{},
				[]piper.DockerVolumeMount{}, false, false, true)
			Expect(err).NotTo(HaveOccurred())

			args := []string{
				"run",
				"--workdir=/tmp/build",
				"--rm",
				"-tty",
				"my-image",
				"my-task.sh",
			}

			Expect(stdout.String()).To(Equal(strings.Join(args, " ") + "\n"))
		})

		It("prints the docker command without running it", func() {
			err := client.Run([]string{"my-task.sh"}, "my-image",
				[]piper.DockerEnv{},
				[]piper.DockerVolumeMount{}, true, true, false)
			Expect(err).NotTo(HaveOccurred())

			args := []string{
				"echo",
				"run",
				"--workdir=/tmp/build",
				"--privileged",
				"-tty",
				"my-image",
				"my-task.sh",
			}

			Expect(stdout.String()).To(Equal(strings.Join(args, " ") + "\n"))
		})

		Context("failure cases", func() {
			Context("when the executable cannot be found", func() {
				It("returns an error", func() {
					client = piper.DockerClient{
						Command: exec.Command("no-such-executable"),
						Stdout:  stdout,
					}
					err := client.Run([]string{"some-command"}, "some-image", []piper.DockerEnv{}, []piper.DockerVolumeMount{}, false, false, false)
					Expect(err).To(MatchError(ContainSubstring("executable file not found in $PATH")))
				})
			})
		})
	})
})
