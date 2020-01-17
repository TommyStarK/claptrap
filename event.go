package main

import "sync"

type event struct {
	mutex sync.Mutex

	name  string
	trace map[uint32]string
}

func (e *event) chmod(timestamp string) {
	e.mutex.Lock()
	e.trace[chmod] = timestamp
	e.mutex.Unlock()
}

func (e *event) create(timestamp string) {
	e.mutex.Lock()
	e.trace[create] = timestamp
	e.mutex.Unlock()
}

func (e *event) remove(timestamp string) {
	e.mutex.Lock()
	e.trace[remove] = timestamp
	e.mutex.Unlock()
}

func (e *event) rename(timestamp string) {
	e.mutex.Lock()
	e.trace[rename] = timestamp
	e.mutex.Unlock()
}

func (e *event) write(timestamp string) {
	e.mutex.Lock()
	e.trace[write] = timestamp
	e.mutex.Unlock()
}

func (e *event) isReadyForProcess() bool {
	e.mutex.Lock()
	timestamp, ok := e.trace[remove]
	e.mutex.Unlock()

	if ok && len(timestamp) > 0 {
		return true
	}

	e.mutex.Lock()
	timestamp, ok = e.trace[rename]
	e.mutex.Unlock()

	if ok && len(timestamp) > 0 {
		return true
	}

	var witness uint32
	e.mutex.Lock()
	for k, v := range e.trace {
		if len(v) > 0 {
			witness = witness | k
		}
	}
	e.mutex.Unlock()

	if witness&write != 0 && witness&chmod != 0 {
		return true
	}

	return false
}
