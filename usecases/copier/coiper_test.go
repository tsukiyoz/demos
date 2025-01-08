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

func TestCopyToPtr(t *testing.T) {
	type Src struct {
		A int
		B string
	}

	type Tgt struct {
		A *int
		B *string
	}

	src := Src{
		A: 1,
		B: "tsukiyo",
	}

	a := 2
	b := "lazywoo"
	copier.Copy(&src, &Tgt{
		A: &a,
		B: &b,
	})

	t.Logf("%v\n", src)
}
