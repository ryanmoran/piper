package piper_test

import (
	"github.com/ryanmoran/piper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VolumeMountBuilder", func() {
	var builder piper.VolumeMountBuilder

	Describe("Build", func() {
		It("builds the volume mounts", func() {
			mounts, err := builder.Build([]string{
				"input-1",
				"input-2",
			}, []string{
				"input-1=/some/path-1",
				"input-2=/some/path-2",
			}, "input")
			Expect(err).NotTo(HaveOccurred())
			Expect(mounts).To(Equal([]piper.DockerVolumeMount{
				{
					LocalPath:  "/some/path-1",
					RemotePath: "/tmp/build/input-1",
				},
				{
					LocalPath:  "/some/path-2",
					RemotePath: "/tmp/build/input-2",
				},
			}))
		})

		Context("failure cases", func() {
			Context("when the input pairs are malformed", func() {
				It("returns an error", func() {
					_, err := builder.Build([]string{}, []string{
						"input-1=something",
						"input-2",
					}, "input")
					Expect(err).To(MatchError("could not parse input \"input-2\". must be of form <input-name>=<input-location>"))
				})
			})

			Context("when an input pair is not specified, but is required", func() {
				It("returns an error", func() {
					_, err := builder.Build([]string{
						"input-1",
						"input-2",
						"input-3",
					}, []string{
						"input-1=/some/path-1",
						"input-3=/some/path-3",
					}, "input")
					Expect(err).To(MatchError("input \"input-2\" is not satisfied. please include an input in command arguments"))
				})
			})
		})
	})
})
