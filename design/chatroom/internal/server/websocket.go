package server

import (
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/tsukiyoz/demos/design/chatroom/internal/logic"
)

func websocketHandleFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("websocket accept failed, err:", err)
		return
	}

	nickname := r.FormValue("nickname")
	if l := len(nickname); l < 2 || l > 20 {
		log.Println("nickname illegal, nickname:", nickname)
		wsjson.Write(r.Context(), conn, "nickname illegal")
		conn.Close(websocket.StatusUnsupportedData, "nickname illegal")
		return
	}

	wsjson.Write(r.Context(), conn, "ok")

	user := logic.NewUser(conn, "token", nickname, r.RemoteAddr)
	logic.Broadcaster.UserEntering(user)

	<-r.Context().Done()

	logic.Broadcaster.UserLeaving(user)
	logic.Broadcaster.Broadcast(&logic.Message{
		UID:      user.ID,
		Type:     logic.MsgTypeUserLeave,
		Content:  nickname + " leave",
		CreateAt: time.Now().UnixMilli(),
	})
	log.Println("user:", user.ID, "leave")
}
