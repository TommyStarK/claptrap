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

	sigchan chan os.Signal
}

func newClaptrap() (*claptrap, error) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGKILL)

	w, err := newWatcher()
	if err != nil {
		return nil, err
	}

	var c = &claptrap{
		events:  make(chan event),
		watcher: w,
		sigchan: sigchan,
	}

	return c, nil
}

func (c *claptrap) add(path string) error {
	return c.watcher.add(path)
}

func (c *claptrap) trap() {
	go c.watcher.watch()

	for {
		select {
		case sig, ok := <-c.sigchan:
			if !ok {

			}

			ch := make(chan struct{})
			c.watcher.stopWatching <- ch
			<-ch
			log.Println("watcher has stopped")
			close(c.events)
			log.Println("closing events chan")
			os.Exit(convertSignalToInt(sig))

		case event, ok := <-c.events:
			if !ok {

			}

			log.Println("claptrap event: ", event)
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
		log.Printf("caught signal: %+v, %s", sig, sig.String())
		rc = 1
	}

	return
}
