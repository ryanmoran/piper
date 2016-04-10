package piper_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPiper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "piper")
}
