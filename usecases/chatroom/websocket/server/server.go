package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "HTTP, hello")
	})

	http.HandleFunc("/ws", handleWebsocket)

	log.Fatal(http.ListenAndServe(":2021", nil))
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "internal error")

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	var v interface{}
	err = wsjson.Read(ctx, conn, &v)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("recv from client: %v\n", v)

	err = wsjson.Write(ctx, conn, "Hello Websocket Client")
	if err != nil {
		log.Println(err)
		return
	}

	conn.Close(websocket.StatusNormalClosure, "")
}
