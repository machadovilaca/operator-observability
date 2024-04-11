package docs

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOperatorDocs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Docs Suite")
}
