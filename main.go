package main

import "log"

func init() {
	log.Println("init")
}

func main() {
	var cfg = &config{
		Path: ".",
	}

	app, err := newClaptrap(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app.trap()
}
