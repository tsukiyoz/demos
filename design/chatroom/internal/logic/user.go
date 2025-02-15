package logic

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
)

type User struct {
	ID       string        `json:"id"`
	NickName string        `json:"nickname"`
	EnterAt  int64         `json:"enter_at"`
	Addr     string        `json:"addr"`
	MsgCh    chan *Message `json:"-"`
	quit     chan error

	conn *websocket.Conn
}

var SystemUser = &User{}

func NewUser(conn *websocket.Conn, token, nickname, addr string) *User {
	user := &User{
		ID:       uuid.New().String(),
		NickName: nickname,
		EnterAt:  time.Now().UnixMilli(),
		Addr:     addr,
		MsgCh:    make(chan *Message, 32),
		conn:     conn,
		quit:     make(chan error, 1),
	}

	return user
}

func (u *User) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			Broadcaster.UserLeaving(u)
			close(u.MsgCh)
			return
		case msg := <-u.MsgCh:
			if msg.Type == MsgTypeNormal {
				msg.Content = u.NickName + ": " + msg.Content
			}
			wsjson.Write(ctx, u.conn, msg)
		}
	}
}

func (u *User) readLoop(ctx context.Context) error {
	var (
		msg map[string]string
		err error
	)
	defer func() {
		u.quit <- err
		close(u.quit)
	}()

	for {
		err = wsjson.Read(ctx, u.conn, &msg)
		if err != nil {
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			} else if errors.Is(err, io.EOF) {
				return nil
			}

			return err
		}

		sendMsg := NewMessage(
			u.ID, MsgTypeNormal,
			msg["content"],
			time.Now().UnixMilli(),
			msg["send_at"],
		)

		Broadcaster.Broadcast(sendMsg)
	}
}

func (u *User) Serve(ctx context.Context) <-chan error {
	go u.readLoop(ctx)
	go u.run(ctx)
	return u.quit
}
