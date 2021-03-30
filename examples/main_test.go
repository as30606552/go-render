package examples

import (
	"os"
	"testing"
)

// Creates directories for output, if there are none.
func TestMain(m *testing.M) {
	if _, err := os.Stat("testdata/pictures"); os.IsNotExist(err) {
		err = os.Mkdir("testdata/pictures", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	if _, err := os.Stat("testdata/output"); os.IsNotExist(err) {
		err = os.Mkdir("testdata/output", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	m.Run()
}
