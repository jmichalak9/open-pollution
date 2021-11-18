package measurement_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMeasurement(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Measurement Suite")
}
