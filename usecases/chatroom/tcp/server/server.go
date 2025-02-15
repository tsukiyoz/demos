package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
)

func main() {
	var (
		lis net.Listener
		err error
	)
	lis, err = net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	go broadcaster()

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConn(conn)
	}
}

type User struct {
	ID      string
	Addr    string
	EnterAt time.Time
	MsgCh   chan string
}

func (u *User) String() string {
	return u.Addr + ", UID:" + u.ID + ", Enter At:" + u.EnterAt.Format(time.DateTime)
}

type Msg struct {
	UID     string
	Content string
}

var (
	registerCh   = make(chan *User)
	unregisterCh = make(chan *User)
	messageCh    = make(chan *Msg)
)

func handleConn(conn net.Conn) {
	defer conn.Close()

	user := User{
		ID:      uuid.New().String(),
		Addr:    conn.RemoteAddr().String(),
		EnterAt: time.Now(),
		MsgCh:   make(chan string, 8),
	}

	go func() {
		for msg := range user.MsgCh {
			fmt.Fprintln(conn, msg)
		}
	}()

	user.MsgCh <- "Welcome, " + user.String()
	userEvtMsg := Msg{
		UID:     user.ID,
		Content: "user:" + user.ID + " has enter",
	}
	registerCh <- &user
	messageCh <- &userEvtMsg

	keepalive := make(chan struct{})
	go func() {
		d := time.Minute
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				conn.Close()
			case <-keepalive:
				timer.Reset(d)
			}
		}
	}()

	var msg Msg
	msg.UID = user.ID
	input := bufio.NewScanner(conn)
	for input.Scan() {
		msg.Content = user.ID + ":" + input.Text()
		messageCh <- &msg
		keepalive <- struct{}{}
	}

	err := input.Err()
	if err != nil {
		log.Printf("read failed, err: %v", err)
	}

	unregisterCh <- &user
	userEvtMsg.Content = "user: " + user.ID + " has left"
	messageCh <- &userEvtMsg
}

func broadcaster() {
	users := make(map[*User]struct{})
	for {
		select {
		case u := <-registerCh:
			users[u] = struct{}{}
		case u := <-unregisterCh:
			delete(users, u)
			close(u.MsgCh)
		case msg := <-messageCh:
			for u := range users {
				if u.ID == msg.UID {
					continue
				}

				u.MsgCh <- msg.Content
			}
		}
	}
}
