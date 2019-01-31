package piper_test

import (
	"io/ioutil"
	"os"

	"github.com/ryanmoran/piper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	Describe("Parse", func() {
		var (
			configFilePath string
			parser         piper.Parser
		)

		BeforeEach(func() {
			tempFile, err := ioutil.TempFile("", "")
			Expect(err).NotTo(HaveOccurred())

			_, err = tempFile.WriteString(`---
image: docker:///some-docker-image
run:
  path: /path/to/run/command
  args: ['-arg1', '-arg2']
inputs:
  - name: input-1
  - name: input-2
  - name: input-3
outputs:
  - name: output-1
  - name: output-2
  - name: output-3
params:
  VAR1: var-1
  VAR2: var-2
`)
			Expect(err).NotTo(HaveOccurred())

			configFilePath = tempFile.Name()

			err = tempFile.Close()
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			err := os.RemoveAll(configFilePath)
			Expect(err).NotTo(HaveOccurred())
		})

		It("parses the task config for the docker image", func() {
			config, err := parser.Parse(configFilePath)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.Image).To(Equal("some-docker-image"))
		})

		It("parses the task config for the run command", func() {
			config, err := parser.Parse(configFilePath)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.Run.Path).To(Equal("/path/to/run/command"))
			Expect(config.Run.Args).To(Equal([]string{"-arg1", "-arg2"}))
		})

		It("parses the task config for the inputs", func() {
			config, err := parser.Parse(configFilePath)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.Inputs).To(Equal([]piper.VolumeMount{
				{Name: "input-1"},
				{Name: "input-2"},
				{Name: "input-3"},
			}))
		})

		It("parses the task config for the outputs", func() {
			config, err := parser.Parse(configFilePath)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.Outputs).To(Equal([]piper.VolumeMount{
				{Name: "output-1"},
				{Name: "output-2"},
				{Name: "output-3"},
			}))
		})

		It("parses the task config for the params", func() {
			config, err := parser.Parse(configFilePath)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.Params).To(Equal(map[string]string{
				"VAR1": "var-1",
				"VAR2": "var-2",
			}))
		})

		It("honors the image_resource", func() {
			tempFile, err := ioutil.TempFile("", "")
			Expect(err).NotTo(HaveOccurred())
			_, err = tempFile.WriteString(`---
image_resource:
  type: docker-image
  source:
    repository: repo/docker-image-name
run:
  path: /path/to/run/command
inputs:
  - name: input-1
  - name: input-2
  - name: input-3
params:
  VAR1: var-1
  VAR2: var-2
`)
			Expect(err).NotTo(HaveOccurred())

			configFilePath = tempFile.Name()

			err = tempFile.Close()
			Expect(err).NotTo(HaveOccurred())
			config, err := parser.Parse(configFilePath)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.Image).To(Equal("repo/docker-image-name"))
		})

		It("honors the image_resource with tags", func() {
			tempFile, err := ioutil.TempFile("", "")
			Expect(err).NotTo(HaveOccurred())
			_, err = tempFile.WriteString(`---
image_resource:
  type: docker-image
  source:
    repository: repo/docker-image-name
    tag: '1.7'
run:
  path: /path/to/run/command
inputs:
  - name: input-1
  - name: input-2
  - name: input-3
params:
  VAR1: var-1
  VAR2: var-2
`)
			Expect(err).NotTo(HaveOccurred())

			configFilePath = tempFile.Name()

			err = tempFile.Close()
			Expect(err).NotTo(HaveOccurred())
			config, err := parser.Parse(configFilePath)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.Image).To(Equal("repo/docker-image-name:1.7"))
		})

		Context("failure cases", func() {
			Context("when the task file does not exist", func() {
				It("returns an error", func() {
					err := os.RemoveAll(configFilePath)
					Expect(err).NotTo(HaveOccurred())

					_, err = parser.Parse(configFilePath)
					Expect(err).To(MatchError(ContainSubstring("no such file or directory")))
				})
			})

			Context("when the task file yaml is not valid", func() {
				It("returns an error", func() {
					err := ioutil.WriteFile(configFilePath, []byte("%%%%%"), 0644)
					Expect(err).NotTo(HaveOccurred())

					_, err = parser.Parse(configFilePath)
					Expect(err).To(MatchError(ContainSubstring("could not find expected directive name")))
				})
			})
		})
	})
})
