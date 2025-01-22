package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	iconn "github.com/tsukiyoz/demos/usecases/websocket/basic/conn"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
var logger = slog.New(slog.NewTextHandler(os.Stderr, nil))

func handler(w http.ResponseWriter, r *http.Request) {
	var (
		conn *iconn.Conn
		err  error
		data []byte
	)
	wsconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	conn, err = iconn.NewIConn(wsconn)
	if err != nil {
		logger.Error("failed to init conn", "err", err)
	}
	defer conn.Close()
	logger.Info("handled a upgraded conn!")

	for {
		data, err = conn.ReadMsg()
		if err != nil {
			if closeErr, ok := err.(*websocket.CloseError); ok {
				logger.Info("websocket closed", "code", closeErr.Code, "Text", closeErr.Text)
				return
			}
		}
		err = conn.SendMsg(data)
		if err != nil {
			logger.Info("write message failed", "err", err)
			return
		}
	}
}

func main() {
	slog.SetDefault(logger)
	logger = logger.With("svc", "websocket")
	http.HandleFunc("/ws", handler)
	http.ListenAndServe(":8080", nil)
}
