package examples

import (
	"computer_graphics/fsutils"
	"testing"
)

// Creates directories for output, if there are none.
func TestMain(m *testing.M) {
	var err = fsutils.MakeDirIfNotExists("testdata/pictures")
	if err != nil {
		panic(err)
	}
	err = fsutils.MakeDirIfNotExists("testdata/output")
	if err != nil {
		panic(err)
	}
	err = fsutils.MakeDirIfNotExists("testdata/pictures/lines")
	if err != nil {
		panic(err)
	}
	err = fsutils.MakeDirIfNotExists("testdata/pictures/triangles")
	if err != nil {
		panic(err)
	}
	err = fsutils.MakeDirIfNotExists("testdata/pictures/vertices")
	if err != nil {
		panic(err)
	}
	err = fsutils.MakeDirIfNotExists("testdata/pictures/imagelibtest")
	if err != nil {
		panic(err)
	}
	m.Run()
}

// Convert the rabbit coordinates to display correctly in the 2000*2000 image.
func defaultRabbitTransformation(x, y, z float64) (float64, float64, float64) {
	return 10000*x + 1000, -10000*y + 1500, 10000*z + 1000
}

// Convert the fox coordinates to display correctly in the 1000*1000 image.
func defaultFoxTransformation(x, y, z float64) (float64, float64, float64) {
	return -5*z + 500, -5*y + 700, 5*x + 500
}
