package parser

import (
	"computer_graphics/fsutils"
	"testing"
)

// Creates directories for output, if there are none.
func TestMain(m *testing.M) {
	if err := fsutils.MakeDirIfNotExists("testdata/output"); err != nil {
		panic(err)
	}
	m.Run()
}
