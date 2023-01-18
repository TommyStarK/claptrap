package main

import (
	"os"
	"testing"

	"github.com/fsnotify/fsnotify"
)

const (
	testContent = `
	this is a content this is a content this is a content this is a content
	this is a content this is a content this is a content this is a content
	this is a content this is a content this is a content this is a content
	this is a content this is a content this is a content this is a content
	this is a content this is a content this is a content this is a content
	this is a content this is a content this is a content this is a content
	this is a content this is a content this is a content this is a content
	this is a content this is a content this is a content this is a content
	this is a content this is a content this is a content this is a content
	`
)

func TestInvalidWatcherInstanciation(t *testing.T) {
	if _, err := newWatcher("", nil, nil); err == nil {
		t.Fatal("watcher instanciation should have failed, invalid path: received empty string")
	}

	if _, err := newWatcher("dummy", nil, nil); err == nil {
		t.Fatal("watcher instanciation should have failed, invalid events channel: received nil")
	}

	eventsChannel := make(chan *event)
	if _, err := newWatcher("dummy", eventsChannel, nil); err == nil {
		t.Fatal("watcher instanciation should have failed, invalid errors channel: received nil")
	}

	errorsChannel := make(chan error)
	if _, err := newWatcher("dummy", eventsChannel, errorsChannel); err == nil {
		t.Fatal("watcher instanciation should have failed, invalid path: does not exist or incorrect permissions")
	}
}

func TestWatcherBehavior(t *testing.T) {
	evch := make(chan *event)
	errch := make(chan error)

	testWatcher, err := newWatcher(testDataPath, evch, errch)
	if err != nil {
		t.Fatal(err)
	}

	triggerWrite := make(chan chan struct{}, 1)
	triggerRename := make(chan chan struct{}, 1)
	triggerRemove := make(chan chan struct{}, 1)
	witness := make(chan struct{}, 1)

	go testWatcher.watch()

	go func(triggerWrite chan chan struct{}) {
		writeDone := <-triggerWrite
		writeFile(testDataPath+"/foo", testContent, errch)
		writeDone <- struct{}{}
	}(triggerWrite)

	go func(triggerRename chan chan struct{}) {
		renameDone := <-triggerRename
		renameFile(testDataPath+"/foo", testDataPath+"/bar", errch)
		renameDone <- struct{}{}
	}(triggerRename)

	go func(triggerRemove chan chan struct{}) {
		removeDone := <-triggerRemove
		removeFile(testDataPath+"/bar", errch)
		removeDone <- struct{}{}
	}(triggerRemove)

	triggerWrite <- witness
	<-witness
	close(triggerWrite)
	processWatcherResult(t, fsnotify.Create, evch)

	triggerRename <- witness
	<-witness
	close(triggerRename)
	processWatcherResult(t, fsnotify.Rename, evch)

	triggerRemove <- witness
	<-witness
	close(triggerRemove)
	processWatcherResult(t, fsnotify.Remove, evch)

	if err := testWatcher.stop(); err != nil {
		t.Fatal(err)
	}

	close(errch)
	close(evch)
	close(witness)
}

func processWatcherResult(t *testing.T, op fsnotify.Op, ch chan *event) {
	event, ok := <-ch

	if !ok || event == nil {
		t.Fatal("unexpected error occurred on event channel")
	}
	if timestamp, ok := event.trace[op]; !ok || len(timestamp) == 0 {
		t.Fatalf("event %s should have been detected and timestamp should not be empty", op.String())
	}
}

func writeFile(path, content string, errchan chan error) {
	var f *os.File
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		errchan <- err
		return
	}

	if err := f.Sync(); err != nil {
		errchan <- err
		return
	}
	if _, err := f.WriteString(content); err != nil {
		errchan <- err
		return
	}
	if err := f.Sync(); err != nil {
		errchan <- err
		return
	}
	if err := f.Close(); err != nil {
		errchan <- err
		return
	}
}

func renameFile(oldpath, newpath string, errchan chan error) {
	if err := os.Rename(oldpath, newpath); err != nil {
		errchan <- err
		return
	}
}

func removeFile(path string, errchan chan error) {
	if err := os.Remove(path); err != nil {
		errchan <- err
		return
	}
}
