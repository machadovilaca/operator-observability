package operatormetrics_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOperatormetrics(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Operatormetrics Suite")
}
