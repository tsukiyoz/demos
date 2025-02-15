package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn)
		log.Println("ok")
		done <- struct{}{}
	}()

	_, err = io.Copy(conn, os.Stdin)
	if err != nil {
		log.Fatal(err)
		return
	}

	conn.Close()
	done <- struct{}{}
}
