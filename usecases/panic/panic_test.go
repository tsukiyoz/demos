package panic

import (
	"testing"
	"time"
)

func TestPanic(t *testing.T) {
	go func() {
		time.Sleep(time.Second)
		t.Logf("triggered panic")
		panic("panic!!!")
	}()

	for {
		select {
		case <-time.After(time.Millisecond * 200):
			t.Logf("main runing...")
		}
	}
}

func TestPanicWithRecover(t *testing.T) {
	handleCrash := func() {
		r := recover()
		if r != nil {
			t.Logf("recover from panic, msg: %v\n", r)
		}
	}

	go func() {
		defer handleCrash()
		time.Sleep(time.Second)
		t.Logf("triggered panic")
		panic("panic!!!")
	}()

	for {
		select {
		case <-time.After(time.Millisecond * 200):
			t.Logf("main runing...")
		}
	}
}
