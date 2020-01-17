package main

import (
	"errors"
	"log"

	"github.com/fsnotify/fsnotify"
)

var eventTypes = map[fsnotify.Op]string{
	fsnotify.Chmod:  "chmod",
	fsnotify.Create: "create",
	fsnotify.Remove: "remove",
	fsnotify.Rename: "rename",
	fsnotify.Write:  "write",
}

type event fsnotify.Event
type watcher struct {
	events     chan event
	fsnWatcher *fsnotify.Watcher

	errs         chan error
	stopWatching chan chan struct{}
}

func newWatcher(path string, events chan event, errs chan error) (*watcher, error) {
	if len(path) == 0 {
		return nil, errors.New("watcher: invalid path, received empty string")
	}

	if events == nil {
		return nil, errors.New("watcher: invalid events channel, received nil")
	}

	if errs == nil {
		return nil, errors.New("watcher: invalid errors channel, received nil")
	}

	fsnWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := fsnWatcher.Add(path); err != nil {
		return nil, err
	}

	var w = &watcher{
		events:       events,
		fsnWatcher:   fsnWatcher,
		errs:         errs,
		stopWatching: make(chan chan struct{}),
	}

	return w, nil
}

func (w *watcher) stop() {
	log.Println("stopping watcher ...")

	close(w.stopWatching)

	if err := w.fsnWatcher.Close(); err != nil {
		log.Printf("failed to close fsnotify watcher: %s", err.Error())
		return
	}

	return
}

func (w *watcher) watch() {
	var notifyWatcherHasStoped chan struct{}

Loop:
	for {
		select {
		case ch := <-w.stopWatching:
			notifyWatcherHasStoped = ch
			break Loop

		case e, ok := <-w.fsnWatcher.Events:
			if !ok {

			}

			w.events <- event(e)

		case err, ok := <-w.fsnWatcher.Errors:
			if !ok {

			}

			w.errs <- err
		}
	}

	w.stop()
	notifyWatcherHasStoped <- struct{}{}
	return
}
