package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type config struct {
	Buttons []button `json:"buttons"`
}

type button struct {
	ID        string `json:"id"`
	ColorCode string `json:"color_code"`
	IconClass string `json:"icon_class"`

	Keycode int `json:"keycode"`
}

func main() {
	configFile := flag.String("config", "config/config.json", "Event config file")
	templateFile := flag.String("template", "template/index.tpl", "Template file")
	outFile := flag.String("outFile", "static/index.html", "Out file")
	flag.Parse()

	b, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var conf config
	err = json.Unmarshal(b, &conf)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	t, err := template.New(path.Base(*templateFile)).ParseFiles(*templateFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	f, err := os.Create(*outFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer func() { _ = f.Close() }()

	err = t.Execute(f, struct {
		Buttons []button
	}{
		Buttons: conf.Buttons,
	})
	if err != nil {
		log.Println(err)
		// TODO: exit with error 1
		return
	}
}
