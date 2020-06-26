package mykind_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMykind(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mykind Suite")
}
