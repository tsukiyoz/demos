package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	SERVER_IP       = "127.0.0.1"
	SERVER_PORT     = 8080
	SERVER_RECV_LEN = 1024
)

func runServer() {
	address := SERVER_IP + ":" + strconv.Itoa(SERVER_PORT)
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	for {
		buf := make([]byte, SERVER_RECV_LEN)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}

		strData := string(buf[:n])
		fmt.Println("Received: ", strData)

		upper := strings.ToUpper(strData)
		_, err = conn.WriteToUDP([]byte(upper), addr)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}

		fmt.Println("Sent: ", upper)
	}
}

func runClient() {
	serverAddr := SERVER_IP + ":" + strconv.Itoa(SERVER_PORT)
	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		panic(err)
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

			_, err = conn.Write([]byte(toWrite))
			if err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}

			buf := make([]byte, SERVER_RECV_LEN)
			n, err = conn.Read(buf)
			if err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}

			fmt.Println("Received: ", string(buf[:n]))
		}
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		runServer()
		wg.Done()
	}()
	go func() {
		runClient()
		wg.Done()
	}()

	wg.Wait()
}
