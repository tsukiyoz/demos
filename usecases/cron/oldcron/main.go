package main

import (
	"fmt"

	oldcron "github.com/robfig/cron"
)

func main() {
	cron := oldcron.New()
	spec := "0 0/1 * * * ?"
	err := cron.AddFunc(spec, func() {
		fmt.Println("hello, world!")
	})
	if err != nil {
		panic(err)
	}
	cron.Run()
}
