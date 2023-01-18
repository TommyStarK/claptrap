package main

import (
	"flag"
	"log"
	"os"
	"plugin"
)

func main() {
	var (
		path       string
		pluginPath string
	)

	flag.StringVar(&path, "path", "", "specify the path to the file/directory to monitor")
	flag.StringVar(&pluginPath, "plugin", "", "path to the plugin to load (.so)")
	flag.Parse()

	if len(path) == 0 {
		log.Fatal("missing file/directory path")
	}

	if len(pluginPath) == 0 {
		log.Fatal("missing plugin path")
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		log.Fatalf("failed to open plugin (%s): %s", pluginPath, err)
	}

	if p == nil {
		log.Fatal("unexpected error occurred, failed to open go plugin")
	}

	handle, err := p.Lookup("Handle")
	if err != nil {
		log.Fatalf("unable to find symbol 'Handle': %s", err)
	}

	if handle == nil {
		log.Fatal("unexpected error occurred, failed to retrieve expected function")
	}

	handler, ok := handle.(func(string, string, string))
	if !ok {
		log.Fatal("unable to retrieve from symbol 'Handle', a function with the following signature: func(string, string, string)")
	}

	app, err := newClaptrap(path, handler)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(app.trap())
}
