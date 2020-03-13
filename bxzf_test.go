package bxzf_test

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/frrad/bxzf"
)

func TestInit(t *testing.T) {
	filename := "test.xz"
	description := fmt.Sprintf("init %s", filename)

	t.Run(description, func(t *testing.T) {
		g := NewGomegaWithT(t)

		var err error
		output := captureOutput(func() {
			_, err = bxzf.OpenFile(filename)
		})

		g.Expect(output).To(Equal(""))
		g.Expect(err).To(BeNil())
	})
}

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}
