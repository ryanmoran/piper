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
			mounts, err := builder.Build([]piper.VolumeMount{
				piper.VolumeMount{Name: "input-1"},
				piper.VolumeMount{Name: "input-2", Optional: true},
				piper.VolumeMount{Name: "output-1"},
				piper.VolumeMount{Name: "output-2"},
				piper.VolumeMount{Path: "cache-1"},
			}, []string{
				"input-1=/some/path-1",
				"input-2=/some/path-2",
			}, []string{
				"output-1=/some/path-3",
				"output-2=/some/path-4",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(mounts[0:4]).To(Equal([]piper.DockerVolumeMount{
				{
					LocalPath:  "/some/path-1",
					RemotePath: "/tmp/build/input-1",
				},
				{
					LocalPath:  "/some/path-2",
					RemotePath: "/tmp/build/input-2",
				},
				{
					LocalPath:  "/some/path-3",
					RemotePath: "/tmp/build/output-1",
				},
				{
					LocalPath:  "/some/path-4",
					RemotePath: "/tmp/build/output-2",
				},
			}))
			Expect(mounts[4].LocalPath).To(Equal("/tmp"))
			Expect(mounts[4].RemotePath).To(Equal("/tmp/build/cache-1"))
		})

		It("expands '~' in paths", func() {
			mounts, err := builder.Build([]piper.VolumeMount{
				piper.VolumeMount{Name: "input-1"},
				piper.VolumeMount{Name: "output-1"},
			}, []string{
				"input-1=~/some/path-1",
			}, []string{
				"output-1=~/some/path-2",
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(mounts[0].LocalPath).ShouldNot(ContainSubstring("~"))
			Expect(mounts[1].LocalPath).ShouldNot(ContainSubstring("~"))
		})

		It("honors the path given in the VolumeMount", func() {
			mounts, err := builder.Build([]piper.VolumeMount{
				piper.VolumeMount{Name: "input-1", Path: "some/path/to/input"},
				piper.VolumeMount{Name: "input-2"},
				piper.VolumeMount{Name: "output-1"},
				piper.VolumeMount{Name: "output-2", Path: "some/path/to/output"},
			}, []string{
				"input-1=/some/path-1",
				"input-2=/some/path-2",
			}, []string{
				"output-1=/some/path-3",
				"output-2=/some/path-4",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(mounts).To(Equal([]piper.DockerVolumeMount{
				{
					LocalPath:  "/some/path-1",
					RemotePath: "/tmp/build/some/path/to/input",
				},
				{
					LocalPath:  "/some/path-2",
					RemotePath: "/tmp/build/input-2",
				},
				{
					LocalPath:  "/some/path-3",
					RemotePath: "/tmp/build/output-1",
				},
				{
					LocalPath:  "/some/path-4",
					RemotePath: "/tmp/build/some/path/to/output",
				},
			}))
		})

		Context("failure cases", func() {
			Context("when the input pairs are malformed", func() {
				It("returns an error", func() {
					_, err := builder.Build([]piper.VolumeMount{}, []string{
						"input-1=something",
						"input-2",
					}, []string{})
					Expect(err).To(MatchError("could not parse input \"input-2\". must be of form <input-name>=<input-location>"))
				})
			})

			Context("when an input pair is not specified, but is required", func() {
				It("returns an error", func() {
					_, err := builder.Build([]piper.VolumeMount{
						{Name: "input-1"},
						{Name: "input-2"},
						{Name: "input-3"},
					}, []string{
						"input-1=/some/path-1",
					}, []string{})
					Expect(err).To(MatchError(`The following required inputs/outputs are not satisfied: input-2, input-3.`))
				})
			})

			Context("when the output pairs are malformed", func() {
				It("returns an error", func() {
					_, err := builder.Build([]piper.VolumeMount{}, []string{}, []string{
						"output-1=something",
						"output-2",
					})
					Expect(err).To(MatchError("could not parse output \"output-2\". must be of form <output-name>=<output-location>"))
				})
			})

			Context("when an input pair is not specified, but is required", func() {
				It("returns an error", func() {
					_, err := builder.Build([]piper.VolumeMount{
						{Name: "output-1"},
						{Name: "output-2"},
						{Name: "output-3"},
					}, []string{}, []string{
						"output-1=/some/path-1",
					})
					Expect(err).To(MatchError(`The following required inputs/outputs are not satisfied: output-2, output-3.`))
				})
			})
		})
	})
})
