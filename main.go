package main

import "log"

func init() {
	log.Println("init")
}

func main() {
	app, err := newClaptrap()
	if err != nil {
		log.Fatal(err)
	}

	app.add(".")
	app.trap()
}
