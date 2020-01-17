package main

import (
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

var evenTypeToString = map[uint32]string{
	create: "create",
	write:  "write",
	remove: "remove",
	rename: "rename",
	chmod:  "chmod",
}

type claptrap struct {
	config  *config
	events  chan *event
	watcher *watcher

	clapMustStop uint32
	errors       chan error
	sigchan      chan os.Signal
	testMode     bool
}

func newClaptrap(cfg *config) (*claptrap, error) {
	errors := make(chan error)
	events := make(chan *event)
	sigchan := make(chan os.Signal, 1)

	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	w, err := newWatcher(cfg.Path, events, errors)
	if err != nil {
		return nil, err
	}

	var c = &claptrap{
		config:  cfg,
		events:  events,
		watcher: w,

		clapMustStop: 0,
		errors:       errors,
		sigchan:      sigchan,
		testMode:     false,
	}

	return c, nil
}

func (c *claptrap) clap(event *event) {
	if atomic.LoadUint32(&c.clapMustStop) == 1 {
		return
	}

	log.Printf("claptrap event on: %s", event.name)
	for k, v := range event.trace {
		log.Printf("-------> Type: %s  @@@  %s", evenTypeToString[k], v)
	}
}

func (c *claptrap) trap() {
	go c.watcher.watch()

	for {
		select {
		case sig, ok := <-c.sigchan:
			if !ok {

			}

			atomic.StoreUint32(&c.clapMustStop, 1)

			if err := c.watcher.stop(); err != nil {
				log.Printf("failed to gracefully stop the watcher: %s", err.Error())
			} else {
				log.Println("watcher has stopped ...")
			}

			close(c.sigchan)
			close(c.errors)
			close(c.events)
			log.Println("claptrap exiting ...")

			if c.testMode {
				return
			}

			os.Exit(convertSignalToInt(sig))

		case event, ok := <-c.events:
			if !ok {

			}

			go c.clap(event)

		case err, ok := <-c.errors:
			if !ok {

			}

			log.Println("error: ", err)
		}
	}
}

func convertSignalToInt(sig os.Signal) (rc int) {
	switch sig.String() {
	case os.Interrupt.String():
		rc = 2 + 128
	case os.Kill.String():
		rc = 9 + 128
	case "terminated":
		rc = 15 + 128
	default:
		rc = 1
	}

	return
}
