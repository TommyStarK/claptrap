package main

import (
	"flag"
	"log"
	"plugin"
)

func main() {

	var (
		path       string
		pluginPath string
	)

	flag.StringVar(&path, "path", "", "specify the path to the file/directory to watch")
	flag.StringVar(&pluginPath, "plugin", "", "path to the plugin to load (.so)")
	flag.Parse()

	if len(path) == 0 {
		log.Fatal("missing file/directory path")
		return
	}

	if len(pluginPath) == 0 {
		log.Fatal("missing plugin path")
		return
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		log.Fatalf("failed to open plugin (%s): %s", pluginPath, err.Error())
		return
	}

	if p == nil {
		log.Fatal("unexpected error occurred")
		return
	}

	handle, err := p.Lookup("Handle")
	if err != nil {
		log.Fatalf("unable to find symbol 'Handle': %s", err.Error())
		return
	}

	if handle == nil {
		log.Fatal("unexpected error occurred")
		return
	}

	handler, ok := handle.(func(string, string, string))
	if !ok {
		log.Fatal("unable to retrieve from symbol func with signature func(string, string, string)")
		return
	}

	app, err := newClaptrap(path, handler)
	if err != nil {
		log.Fatal(err)
		return
	}

	app.trap()
	return
}
