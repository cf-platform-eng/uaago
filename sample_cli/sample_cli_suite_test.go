package sample_cli_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSampleCli(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SampleCli Suite")
}
