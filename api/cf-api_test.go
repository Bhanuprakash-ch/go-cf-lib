package api

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestCfApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cf api Suite")
}
