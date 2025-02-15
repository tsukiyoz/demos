package logic

import "strconv"

type MsgType int

const (
	MsgTypeNormal MsgType = iota
	MsgTypeWelcome
	MsgTypeUserEnter
	MsgTypeUserLeave
	MsgTypeError
)

type Message struct {
	UID      string   `json:"uid"`
	Type     MsgType  `json:"type"`
	Content  string   `json:"content"`
	CreateAt int64    `json:"create_at"`
	SendAt   int64    `json:"send_at"`
	Ats      []string `json:"ats"`
}

func NewMessage(uid string, typ MsgType, content string, createAt int64, sendAt string) *Message {
	sendAtTime, _ := strconv.ParseInt(sendAt, 10, 64)
	return &Message{
		UID:      uid,
		Type:     typ,
		Content:  content,
		CreateAt: createAt,
		SendAt:   sendAtTime,
	}
}
