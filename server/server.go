package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/micmonay/keybd_event"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Server .
type Server struct {
	configFile string
	keyboard keybd_event.KeyBonding

	keyMapLock sync.Mutex
	keyMap map[string]int
}

// New .
func New(configFile string) (*Server, error) {
	keyboard, err := keybd_event.NewKeyBonding()
	if err != nil {
		return nil, err
	}

	s := &Server{
		configFile: configFile,
		keyboard: keyboard,
		keyMap: make(map[string]int),
	}

	return s, s.LoadConfig()
}

type config struct {
	Events map[string]int `json:"events"`
}

func (s *Server) LoadConfig() error {
	b, err := ioutil.ReadFile(s.configFile)
	if err != nil {
		return err
	}

	var conf config
	err = json.Unmarshal(b, &conf)
	if err != nil {
		return err
	}

	s.keyMapLock.Lock()
	s.keyMap = conf.Events
	s.keyMapLock.Unlock()

	return nil
}

func (s *Server) Routes() *http.ServeMux{
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

		s.keyMapLock.Lock()
		key, ok := s.keyMap[string(message)]
		s.keyMapLock.Unlock()

		if !ok {
			log.Println("invalid payload", string(message))
			continue
		}
		log.Printf("recv keycode: %d, Message: %s", key, message)

		s.keyboard.SetKeys(key)
		log.Println("Keypress", key)
		err = s.keyboard.Launching()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Keypress done", key)
	}
}