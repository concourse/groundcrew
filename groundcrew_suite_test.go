package groundcrew_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGroundcrew(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Groundcrew Suite")
}
