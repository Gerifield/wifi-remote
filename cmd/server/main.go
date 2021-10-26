package main

import (
	"flag"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
	"github.com/micmonay/keybd_event"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var keyMap = map[string]int{
	"1": keybd_event.VK_F13,
	//"2": keybd_event.VK_B,
}

func connect(w http.ResponseWriter, r *http.Request) {
	log.Println("New client")
	c, err := upgrader.Upgrade(w, r, nil)
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

		key, ok := keyMap[string(message)]
		if !ok {
			continue
		}
		log.Printf("recv keycode: %d, Message: %s", key, message)

		keyboard.SetKeys(key)
		log.Println("Keypress", key)
		err = keyboard.Launching()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Keypress done", key)
	}
}

var keyboard keybd_event.KeyBonding

func init() {
	var err error
	keyboard, err = keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}
}

func main() {
	addr := flag.String("listen", ":8080", "HTTP listen endpoint")
	flag.Parse()

	http.HandleFunc("/connect", connect)
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello"))
	})

	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	log.Println("Started")
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Println(err)
	}
}
