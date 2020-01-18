package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

type claptrap struct {
	events  chan *event
	handler func(string, string, string)
	watcher *watcher

	clapMustStop uint32
	errors       chan error
	sigchan      chan os.Signal
	target       string
	testMode     bool
}

func newClaptrap(path string, handler func(string, string, string)) (*claptrap, error) {
	errors := make(chan error)
	events := make(chan *event)
	sigchan := make(chan os.Signal, 1)

	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	w, err := newWatcher(path, events, errors)
	if err != nil {
		return nil, err
	}

	var c = &claptrap{
		events:  events,
		handler: handler,
		watcher: w,

		clapMustStop: 0,
		errors:       errors,
		sigchan:      sigchan,
		target:       path,
		testMode:     false,
	}

	return c, nil
}

func (c *claptrap) clap(event *event) {
	if atomic.LoadUint32(&c.clapMustStop) == 1 {
		return
	}

	if c.testMode {
		return
	}

	var (
		action    = "UPDATE"
		timestamp = ""
	)

	for fsnotifyEventType, ts := range event.trace {
		switch fsnotifyEventType.String() {
		case fsnotify.Create.String():
			action = fsnotify.Create.String()
			timestamp = ts
		case fsnotify.Remove.String():
			action = fsnotify.Remove.String()
			timestamp = ts
		case fsnotify.Rename.String():
			action = fsnotify.Rename.String()
			timestamp = ts
		}
	}

	if len(timestamp) == 0 {
		if ts, ok := event.trace[fsnotify.Chmod]; ok && len(ts) > 0 {
			timestamp = ts
		}
	}

	go c.handler(action, event.name, timestamp)
	return
}

func (c *claptrap) trap() {
	go c.watcher.watch()

	for {
		select {
		case sig, ok := <-c.sigchan:
			if !ok {
				os.Exit(1)
				return
			}

			atomic.StoreUint32(&c.clapMustStop, 1)

			if err := c.watcher.stop(); err != nil {
				log.Printf("failed to gracefully stop the watcher: %s", err.Error())
			}

			close(c.sigchan)
			close(c.errors)
			close(c.events)

			if c.testMode {
				return
			}

			os.Exit(convertSignalToInt(sig))
			return

		case event, ok := <-c.events:
			if ok && event != nil {
				go c.clap(event)
			}

		case err, ok := <-c.errors:
			if ok && err != nil {
				log.Println("error: ", err)
			}
		}
	}
}

func convertSignalToInt(sig os.Signal) (rc int) {
	rc = 1

	if sig == nil {
		return
	}

	switch sig.String() {
	case os.Interrupt.String():
		rc = 2
	case os.Kill.String():
		rc = 9
	case "terminated":
		rc = 15
	default:
	}

	return
}
