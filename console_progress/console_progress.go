package main

import (
	_ "embed"
	"github.com/cheggaaa/pb/v3"
	"math/rand"
	"time"
)

//go:embed tmpl.txt
var tmp string

func main() {
	count := 500
	bar := pb.ProgressBarTemplate(tmp).Start(count)
	for i := 0; i < count; i++ {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(20)))
		bar.Increment()
	}
	bar.Finish()
}
