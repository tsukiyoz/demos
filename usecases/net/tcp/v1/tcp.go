package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	SERVER_IP       = "127.0.0.1"
	SERVER_PORT     = 8080
	SERVER_RECV_LEN = 1024
)

func runServer() {
	address := SERVER_IP + ":" + strconv.Itoa(SERVER_PORT)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer lis.Close()

	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		defer conn.Close()

		for {
			data := make([]byte, SERVER_RECV_LEN)
			_, err = conn.Read(data)
			if err != nil {
				fmt.Println(err)
				break
			}

			strData := string(data)
			fmt.Println("Received: ", strData)

			upper := strings.ToUpper(strData)
			_, err = conn.Write([]byte(upper))
			if err != nil {
				fmt.Println(err)
				break
			}

			fmt.Println("Send: ", upper)
		}
	}
}

func runClient() {
	addr := SERVER_IP + ":" + strconv.Itoa(SERVER_PORT)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer conn.Close()

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()

		n := 0
		for w := 0; w < len(line); w += n {
			var toWrite string
			if len(line)-w > SERVER_RECV_LEN {
				toWrite = line[w : w+SERVER_RECV_LEN]
			} else {
				toWrite = line[w:]
			}

			n, err = conn.Write([]byte(toWrite))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("Write: ", toWrite)

			buf := make([]byte, SERVER_RECV_LEN)
			_, err = conn.Read(buf)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("Recv: ", string(buf))
		}
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		runServer()
	}()
	time.Sleep(time.Second)
	go func() {
		defer wg.Done()
		runClient()
	}()

	wg.Wait()
}
