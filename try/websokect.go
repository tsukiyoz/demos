package main

import "github.com/gorilla/websocket"

func main() {
	dialer := websocket.Dialer{}
	dialer.NetDialContext()
}
