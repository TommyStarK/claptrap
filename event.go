package main

import (
	"sync"

	"github.com/fsnotify/fsnotify"
)

type event struct {
	mutex sync.Mutex

	name  string
	trace map[fsnotify.Op]string
}

func (e *event) chmod(timestamp string) {
	e.mutex.Lock()
	e.trace[fsnotify.Chmod] = timestamp
	e.mutex.Unlock()
}

func (e *event) create(timestamp string) {
	e.mutex.Lock()
	e.trace[fsnotify.Create] = timestamp
	e.mutex.Unlock()
}

func (e *event) remove(timestamp string) {
	e.mutex.Lock()
	e.trace[fsnotify.Remove] = timestamp
	e.mutex.Unlock()
}

func (e *event) rename(timestamp string) {
	e.mutex.Lock()
	e.trace[fsnotify.Rename] = timestamp
	e.mutex.Unlock()
}

func (e *event) write(timestamp string) {
	e.mutex.Lock()
	e.trace[fsnotify.Write] = timestamp
	e.mutex.Unlock()
}

func (e *event) isReadyForBeingProcessed() bool {
	e.mutex.Lock()
	timestamp, ok := e.trace[fsnotify.Remove]
	e.mutex.Unlock()

	if ok && len(timestamp) > 0 {
		return true
	}

	e.mutex.Lock()
	timestamp, ok = e.trace[fsnotify.Rename]
	e.mutex.Unlock()

	if ok && len(timestamp) > 0 {
		return true
	}

	var witness fsnotify.Op
	e.mutex.Lock()
	for k, v := range e.trace {
		if len(v) > 0 {
			witness = witness | k
		}
	}
	e.mutex.Unlock()

	if witness&fsnotify.Write != 0 && witness&fsnotify.Chmod != 0 {
		return true
	}

	return false
}
