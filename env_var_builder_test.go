package piper_test

import (
	"github.com/ryanmoran/piper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EnvVarBuilder", func() {
	It("returns a list of environment variables", func() {
		vars := piper.EnvVarBuilder{}.Build([]string{
			"VAR1=var-1",
			"VAR3=var-3",
		}, map[string]string{
			"VAR1": "default-var-1",
			"VAR2": "default-var-2",
		})
		Expect(vars).To(ConsistOf([]piper.DockerEnv{
			{
				Key:   "VAR1",
				Value: "var-1",
			},
			{
				Key:   "VAR2",
				Value: "default-var-2",
			},
		}))

	})

	Context("when env vars have '=' signs in the value", func() {
		It("returns a list of environment variables with '=' signs still in their place", func() {
			vars := piper.EnvVarBuilder{}.Build([]string{
				"VAR1=var-1==42=",
			}, map[string]string{
				"VAR1": "meow",
			})
			Expect(vars).To(ConsistOf([]piper.DockerEnv{
				{
					Key:   "VAR1",
					Value: "var-1==42=",
				},
			}))
		})
	})
})
