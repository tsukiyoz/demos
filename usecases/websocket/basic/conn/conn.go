package iconn

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Conn struct {
	ws        *websocket.Conn
	inCh      chan []byte
	outCh     chan []byte
	close     chan struct{}
	closeOnce sync.Once
}

func NewIConn(ws *websocket.Conn) (*Conn, error) {
	c := &Conn{
		ws:    ws,
		inCh:  make(chan []byte, 1024),
		outCh: make(chan []byte, 1024),
		close: make(chan struct{}, 1),
	}
	go c.readLoop()
	go c.writeLoop()
	return c, nil
}

func (c *Conn) readLoop() {
	defer c.Close()
	var data []byte
	var err error
	for {
		_, data, err = c.ws.ReadMessage()
		if err != nil {
			if _, ok := <-c.close; ok {
				log.Printf("Conn is closed, stop read message\n")
				return
			}
			log.Fatalf("failed to read message, %v", err)
			return
		}
		select {
		case <-c.close:
			log.Printf("Conn is closed, stop read message\n")
			return
		case c.inCh <- data:
		}
	}
}

func (c *Conn) writeLoop() {
	defer c.Close()
	var data []byte
	var err error
	for {
		select {
		case <-c.close:
			log.Printf("Conn is closed, stop write message\n")
			return
		case data = <-c.outCh:
		}
		err = c.ws.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			if _, ok := <-c.close; ok {
				log.Printf("Conn is closed, stop write message\n")
				return
			}
			log.Fatalf("failed to write message, %v", err)
			return
		}
	}
}

// ----------------- API -------------------

func (c *Conn) ReadMsg() (data []byte, err error) {
	data = <-c.inCh
	return data, nil
}

func (c *Conn) SendMsg(data []byte) (err error) {
	c.outCh <- data
	return nil
}

func (c *Conn) Close() error {
	c.closeOnce.Do(func() {
		close(c.close)
		c.ws.Close()
	})
	return nil
}
