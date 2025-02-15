package logic

import (
	"log"
)

var Broadcaster = NewBroadcaster()

type broadcaster struct {
	users map[string]*User

	enterCh chan *User
	leaveCh chan *User
	msgCh   chan *Message

	// checkUserCh      chan string
	// checkUserCanInCh chan bool
}

func NewBroadcaster() *broadcaster {
	b := &broadcaster{
		users:   make(map[string]*User),
		enterCh: make(chan *User),
		leaveCh: make(chan *User),
		msgCh:   make(chan *Message, 16),
	}

	go b.run()

	return b
}

func (b *broadcaster) run() {
	for {
		select {
		case u := <-b.enterCh:
			b.users[u.NickName] = u
		case u := <-b.leaveCh:
			delete(b.users, u.NickName)
			close(u.MsgCh)
		case msg := <-b.msgCh:
			for _, u := range b.users {
				if u.ID == msg.UID {
					continue
				}
				u.MsgCh <- msg
			}
		}
	}
}

func (b *broadcaster) UserEntering(u *User) {
	b.enterCh <- u
}

func (b *broadcaster) UserLeaving(u *User) {
	b.leaveCh <- u
}

func (b *broadcaster) Broadcast(msg *Message) {
	if len(b.msgCh) >= 1024 {
		log.Println("broadcast queue out of limit")
	}
	b.msgCh <- msg
}
