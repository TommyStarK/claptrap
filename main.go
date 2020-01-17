package main

import "log"

func init() {
	log.Println("init")
}

func main() {
	var cfg = &config{
		path: ".",
	}

	app, err := newClaptrap(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app.trap()
}
