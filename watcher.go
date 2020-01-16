package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type watcher struct {
	fsnWatcher *fsnotify.Watcher
	events     chan event

	stopWatching chan chan struct{}
}

func newWatcher() (*watcher, error) {
	fsnWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	var w = &watcher{
		fsnWatcher:   fsnWatcher,
		stopWatching: make(chan chan struct{}),
	}

	return w, nil
}

func (w *watcher) add(path string) error {
	return w.fsnWatcher.Add(path)
}

func (w *watcher) stop() {
	if err := w.fsnWatcher.Close(); err != nil {
		log.Printf("failed to close fsnotify watcher: %s", err.Error())
	}
}

func (w *watcher) watch() {
	var notifyWatcherHasStoped chan struct{}

	for {
		select {
		case ch := <-w.stopWatching:
			notifyWatcherHasStoped = ch
			goto stop

		case event := <-w.fsnWatcher.Events:
			log.Println("event: ", event)

		case err := <-w.fsnWatcher.Errors:
			log.Println("error: ", err)
		}
	}

stop:
	log.Println("stopping watcher")
	w.stop()
	notifyWatcherHasStoped <- struct{}{}
}
