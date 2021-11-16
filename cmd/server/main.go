package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gerifield/wifi-remote/remote"
	"github.com/gerifield/wifi-remote/server"
	"github.com/micmonay/keybd_event"
)

func main() {
	addr := flag.String("listen", ":8080", "HTTP listen endpoint")
	configFile := flag.String("config", "config/config.json", "Event config file")
	flag.Parse()

	keyboard, err := keybd_event.NewKeyBonding()
	if err != nil {
		log.Println(err)
		return
	}

	r, err := remote.New(*configFile, &keyboard)
	if err != nil {
		log.Println(err)
		return
	}

	srv := server.New(r)

	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	log.Println("Started", *addr)
	ch := make(chan os.Signal, 1)
	go func() {
		for range ch {
			log.Println("config reload")
			err := r.LoadConfig()
			if err != nil {
				log.Println(err)
			}
		}
	}()
	signal.Notify(ch, syscall.SIGUSR1)

	err = http.ListenAndServe(*addr, srv.Routes())
	if err != nil {
		log.Println(err)
	}
}
