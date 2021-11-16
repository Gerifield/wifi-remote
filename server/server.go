package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type keyPresser interface {
	KeyPress(button string) error
}

// Server .
type Server struct {
	remote keyPresser
}

// New .
func New(remote keyPresser) *Server {
	return &Server{
		remote: remote,
	}
}

func (s *Server) Routes() *http.ServeMux {
	m := http.NewServeMux()

	m.HandleFunc("/connect", s.handleConnect)
	m.HandleFunc("/", s.handleRoot)

	return m
}

func (s *Server) handleRoot(rw http.ResponseWriter, _ *http.Request) {
	rw.Write([]byte("Hello"))
}

func (s *Server) handleConnect(rw http.ResponseWriter, r *http.Request) {
	log.Println("New client")
	c, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read err:", err)
			break
		}

		log.Println("received button:", string(message))
		err = s.remote.KeyPress(string(message))
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("keypress done")
	}
}
