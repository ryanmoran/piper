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

			Expect(config.Command).To(Equal("/path/to/run/command"))
		})

		It("parses the task config for the inputs", func() {
			config, err := parser.Parse(configFilePath)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.Inputs).To(Equal([]string{
				"input-1",
				"input-2",
				"input-3",
			}))
		})

		It("parses the task config for outputs", func(){
			config, err := parser.Parse(configFilePath)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.Outputs).To(Equal([]string{
				"output-1",
				"output-2",
				"output-3",
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
