package main

import (
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

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

	w, err := newWatcher(cfg.path, events, errors)
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
	for fsnotifyEventType, timestamp := range event.trace {
		log.Printf("-------> Type: %s  @@@  %s", fsnotifyEventType.String(), timestamp)
	}
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

			log.Println("claptrap exiting ...")

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
