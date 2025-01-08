package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestCaptureKillSignal(t *testing.T) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigchan:
		t.Logf("嘻嘻我收到了一条终止指令\n")
	}
	for i := 0; i < 5; i++ {
		t.Logf("我正在退出...\n")
		time.Sleep(time.Second)
	}
	fmt.Println("我退出了！")
}

func TestCaptureKillSignalForever(t *testing.T) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-sigchan:
			t.Logf("嘻嘻我收到了一条终止指令\n")
		}
	}
}
