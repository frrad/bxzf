package bxzf_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/frrad/bxzf"
)

func TestInit(t *testing.T) {
	filename := "test.xz"
	description := fmt.Sprintf("init %s", filename)

	t.Run(description, func(t *testing.T) {
		g := NewGomegaWithT(t)

		_, err := bxzf.OpenFile(filename)

		g.Expect(err).To(BeNil())
	})
}
