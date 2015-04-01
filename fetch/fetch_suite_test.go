package fetch_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFetch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fetch Suite")
}
