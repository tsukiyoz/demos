package copier

import (
	"testing"

	"github.com/jinzhu/copier"
)

type Point struct {
	X, Y int
}

type SRC struct {
	Name string
	Point
}

type DST struct {
	Val struct {
		Name string
		X    int
	}
	Y int
}

func TestCopy(t *testing.T) {
	src := SRC{
		Name: "foo",
		Point: Point{
			X: 1,
			Y: 2,
		},
	}
	dst := DST{}
	copier.Copy(&dst, &src)
	t.Logf("%+v", dst)
}
