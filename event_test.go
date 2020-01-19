package main

import (
	"sync"
	"testing"

	"github.com/fsnotify/fsnotify"
)

func TestEventIsReadyForBeingProcessed(t *testing.T) {
	var event1 = &event{
		mutex: sync.Mutex{},
		name:  "dummy",
		trace: make(map[fsnotify.Op]string),
	}

	event1.create("fake-timestamp-1")
	event1.write("fake-timestamp-2")

	if event1.isReadyForBeingProcessed() {
		t.Log("event1 should not be ready for being processed, missing 'CHMOD' event, meaning the file is still being edited")
		t.Fail()
	}

	event1.chmod("fake-timestamp-3")
	if !event1.isReadyForBeingProcessed() {
		t.Log("event1 should be ready for being processed")
		t.Fail()
	}

	var event2 = &event{
		mutex: sync.Mutex{},
		name:  "dummy2",
		trace: make(map[fsnotify.Op]string),
	}

	event2.remove("fake-timestamp-1")
	if !event2.isReadyForBeingProcessed() {
		t.Log("event2 should be ready for being processed")
		t.Fail()
	}

	var event3 = &event{
		mutex: sync.Mutex{},
		name:  "dummy3",
		trace: make(map[fsnotify.Op]string),
	}

	event3.rename("fake-timestamp-1")
	if !event3.isReadyForBeingProcessed() {
		t.Log("event3 should be ready for being processed")
		t.Fail()
	}
}
