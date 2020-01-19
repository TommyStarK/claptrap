package main

import (
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

type claptrap struct {
	events  chan *event
	handler func(string, string, string)
	watcher *watcher

	clapMustStop uint32
	errors       chan error
	sigchan      chan os.Signal
	target       string

	// for sake of tests
	testchan chan [3]string
	testMode bool
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
		// for sake of tests
		testchan: nil,
		testMode: false,
	}

	return c, nil
}

func (c *claptrap) clap(event *event) {
	var (
		action    = "UPDATE"
		timestamp = ""
	)

	for fsnotifyEventType, ts := range event.trace {
		switch fsnotifyEventType.String() {
		case fsnotify.Remove.String():
			action = fsnotify.Remove.String()
			timestamp = ts
			break
		case fsnotify.Rename.String():
			action = fsnotify.Rename.String()
			timestamp = ts
			break
		}
	}

	if len(timestamp) == 0 {
		if ts, ok := event.trace[fsnotify.Chmod]; ok && len(ts) > 0 {
			timestamp = ts
		}

		if ts, ok := event.trace[fsnotify.Create]; ok && len(ts) > 0 {
			action = fsnotify.Create.String()
		}
	}

	if !c.testMode && atomic.LoadUint32(&c.clapMustStop) == 0 {
		go c.handler(action, event.name, timestamp)
		return
	}

	// for sake of tests
	if c.testchan != nil && atomic.LoadUint32(&c.clapMustStop) == 0 {
		c.testchan <- [3]string{action, event.name, timestamp}
	}

	return
}

func (c *claptrap) trap() {
	go c.watcher.watch()

	for {
		select {
		case sig, ok := <-c.sigchan:
			if !ok {
				log.Fatal("unexpected error occurred on signal channel")
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
				log.Printf("error: %s", err.Error())
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
