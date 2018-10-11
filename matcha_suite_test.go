package matcha_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMatcha(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Matcha Suite")
}
