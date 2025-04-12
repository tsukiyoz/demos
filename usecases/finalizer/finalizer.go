package main

import (
	"fmt"
	"runtime"
	"time"
)

type App struct {
	closed chan struct{}
}

func New() *App {
	app := &App{
		closed: make(chan struct{}),
	}
	closed := app.closed
	runtime.SetFinalizer(app, func(a *App) {
		fmt.Print("Finalizer called\n")
		close(closed)
	})
	go doLoop(closed)
	return app
}

func (a *App) Version() string {
	return "1.0.0"
}

func doLoop(closeCh <-chan struct{}) {
	ticker := time.NewTicker(400 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			fmt.Print(".")
		case <-closeCh:
			ticker.Stop()
			fmt.Print("Finalizer loop stopped\n")
			return
		}
	}
}

func main() {
	app := New()
	// Simulate some work
	time.Sleep(2 * time.Second)
	fmt.Println("Version:", app.Version())
	app = nil // remove reference to app
	runtime.GC()
	time.Sleep(5 * time.Second)
}
