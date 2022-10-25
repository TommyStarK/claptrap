package main

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
)

type watcher struct {
	fsnWatcher *fsnotify.Watcher

	errs         chan error
	events       chan *event
	rwmutex      sync.RWMutex
	stopWatching chan chan struct{}
	trace        map[string]*event

	processingMustStop uint32
}

func newWatcher(path string, events chan *event, errs chan error) (*watcher, error) {
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
		defer fsnWatcher.Close()
		return nil, err
	}

	var w = &watcher{
		fsnWatcher: fsnWatcher,

		errs:         errs,
		events:       events,
		rwmutex:      sync.RWMutex{},
		stopWatching: make(chan chan struct{}),
		trace:        make(map[string]*event),

		processingMustStop: 0,
	}

	return w, nil
}

func (w *watcher) processEvent(fsevent fsnotify.Event) {
	var (
		exist       bool
		targetEvent *event
		timestamp   string = time.Now().UTC().String()
	)

	w.rwmutex.RLock()
	targetEvent, exist = w.trace[fsevent.Name]
	w.rwmutex.RUnlock()

	if !exist {
		var newEvent = &event{
			mutex: sync.Mutex{},
			trace: make(map[fsnotify.Op]string),

			name: fsevent.Name,
		}

		targetEvent = newEvent
		w.rwmutex.Lock()
		w.trace[fsevent.Name] = targetEvent
		w.rwmutex.Unlock()
	}

	switch fsevent.Op {
	case fsnotify.Create:
		targetEvent.create(timestamp)
	case fsnotify.Write:
		targetEvent.write(timestamp)
	case fsnotify.Remove:
		targetEvent.remove(timestamp)
	case fsnotify.Rename:
		targetEvent.rename(timestamp)
	case fsnotify.Chmod:
		targetEvent.chmod(timestamp)
	}

	if targetEvent.isReadyForBeingProcessed() && atomic.LoadUint32(&w.processingMustStop) == 0 {
		w.rwmutex.Lock()
		w.events <- targetEvent
		delete(w.trace, targetEvent.name)
		w.rwmutex.Unlock()
	}

	return
}

func (w *watcher) stop() error {
	atomic.StoreUint32(&w.processingMustStop, 1)
	ch := make(chan struct{})
	w.stopWatching <- ch
	<-ch
	close(ch)
	close(w.stopWatching)
	return w.fsnWatcher.Close()
}

func (w *watcher) watch() {
	for {
		select {
		case ch := <-w.stopWatching:
			ch <- struct{}{}
			return
		case e, ok := <-w.fsnWatcher.Events:
			if ok {
				go w.processEvent(e)
			}
		case err, ok := <-w.fsnWatcher.Errors:
			if ok && err != nil {
				w.errs <- err
			}
		}
	}
}
