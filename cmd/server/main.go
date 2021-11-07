package main

import (
	"flag"
	"github.com/gerifield/wifi-remote/server"
	"log"
	"net/http"
	"runtime"
	"time"
)

func main() {
	addr := flag.String("listen", ":8080", "HTTP listen endpoint")
	configFile := flag.String("config", "config.json", "Event config file")
	flag.Parse()

	srv, err := server.New(*configFile)
	if err != nil {
		log.Println(err)
		return
	}

	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	log.Println("Started", *addr)
	err = http.ListenAndServe(*addr, srv.Routes())
	if err != nil {
		log.Println(err)
	}
}
