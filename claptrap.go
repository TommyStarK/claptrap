package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type claptrap struct {
	events  chan event
	watcher *watcher

	errors   chan error
	sigchan  chan os.Signal
	testMode bool
}

func newClaptrap(cfg *config) (*claptrap, error) {
	errors := make(chan error)
	events := make(chan event)
	sigchan := make(chan os.Signal, 1)

	signal.Notify(sigchan,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGKILL)

	w, err := newWatcher(cfg.Path, events, errors)
	if err != nil {
		return nil, err
	}

	var c = &claptrap{
		events:   events,
		watcher:  w,
		errors:   errors,
		sigchan:  sigchan,
		testMode: false,
	}

	return c, nil
}

func (c *claptrap) trap() {
	var signalCaught os.Signal

	go c.watcher.watch()

Loop:
	for {
		select {
		case sig, ok := <-c.sigchan:
			if !ok {

			}

			signalCaught = sig
			break Loop

		case event, ok := <-c.events:
			if !ok {

			}

			eventType, ok := eventTypes[event.Op]
			if !ok {

			}

			log.Printf("claptrap event %s - file: %s", eventType, event.Name)

		case err, ok := <-c.errors:
			if !ok {

			}

			log.Println("error: ", err)
		}
	}

	ch := make(chan struct{})
	c.watcher.stopWatching <- ch
	<-ch
	close(ch)
	close(c.sigchan)
	close(c.errors)
	close(c.events)
	log.Println("claptrap exiting ...")

	if c.testMode {
		return
	}

	os.Exit(convertSignalToInt(signalCaught))
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
		// log.Printf("caught signal: %+v, %s", sig, sig.String())
		rc = 1
	}

	return
}
