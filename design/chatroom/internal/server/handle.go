package server

import (
	_ "embed"
	"net/http"
)

//go:embed home.html
var homeTemplate string

type Server struct{}

func NewServer() *Server {
	srv := &Server{}
	return srv
}

func (s *Server) RegisterHandle() {
	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", websocketHandleFunc)
}
