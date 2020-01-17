package main

import (
	"testing"
)

func TestInvalidWatcherInstanciation(t *testing.T) {
	if _, err := newWatcher("", nil, nil); err == nil {
		t.Log("watcher instanciation should have failed, invalid path: received empty string")
		t.Fail()
	}

	if _, err := newWatcher("dummy", nil, nil); err == nil {
		t.Log("watcher instanciation should have failed, invalid events channel: received nil")
		t.Fail()
	}

	eventsChannel := make(chan *event)
	if _, err := newWatcher("dummy", eventsChannel, nil); err == nil {
		t.Log("watcher instanciation should have failed, invalid errors channel: received nil")
		t.Fail()
	}

	errorsChannel := make(chan error)
	if _, err := newWatcher("dummy", eventsChannel, errorsChannel); err == nil {
		t.Log("watcher instanciation should have failed, invalid path: does not exist or incorrect permissions")
		t.Fail()
	}
}
