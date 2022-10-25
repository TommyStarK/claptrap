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

func TestWatcherBehavior(t *testing.T) {
	evch := make(chan *event)
	errch := make(chan error)

	testWatcher, err := newWatcher(testDataPath, evch, errch)
	if err != nil {
		t.Fatal(err)
	}

	stopWatchErr := make(chan chan struct{})
	go func() {
		for {
			select {
			case ch := <-stopWatchErr:
				ch <- struct{}{}
				return
			case err, ok := <-errch:
				if !ok || err == nil {
					t.Log("unexpected error occurred on error channel")
					t.Fail()
					continue
				}
				t.Log(err)
				t.Fail()
			}
		}
	}()

	go testWatcher.watch()

	triggerWrite := make(chan chan struct{})
	go func() {
		writeDone := <-triggerWrite
		writeFile(testDataPath+"/foo", testContent, errch)
		writeDone <- struct{}{}
		return
	}()

	triggerRename := make(chan chan struct{})
	go func() {
		renameDone := <-triggerRename
		renameFile(testDataPath+"/foo", testDataPath+"/bar", errch)
		renameDone <- struct{}{}
	}()

	triggerRemove := make(chan chan struct{})
	go func() {
		removeDone := <-triggerRemove
		removeFile(testDataPath+"/bar", errch)
		removeDone <- struct{}{}
		return
	}()

	witness := make(chan struct{})
	triggerWrite <- witness
	<-witness
	close(triggerWrite)
	processWatcherResult(fsnotify.Create, evch, t)

	triggerRename <- witness
	<-witness
	close(triggerRename)
	processWatcherResult(fsnotify.Rename, evch, t)

	triggerRemove <- witness
	<-witness
	close(triggerRemove)
	processWatcherResult(fsnotify.Remove, evch, t)

	stopWatchErr <- witness
	<-witness
	close(stopWatchErr)

	if err := testWatcher.stop(); err != nil {
		t.Log(err)
		t.Fail()
	}

	close(errch)
	close(evch)
	close(witness)
}

func processWatcherResult(op fsnotify.Op, ch chan *event, t *testing.T) {
	event, ok := <-ch

	if !ok || event == nil {
		t.Log("unexpected error occurred on event channel")
		t.Fail()
		return
	}
	if timestamp, ok := event.trace[op]; !ok || len(timestamp) == 0 {
		t.Logf("event %s should have been detected and timestamp should not be empty", op.String())
		t.Fail()
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
	return
}

func renameFile(oldpath, newpath string, errchan chan error) {
	if err := os.Rename(oldpath, newpath); err != nil {
		errchan <- err
		return
	}
	return
}

func removeFile(path string, errchan chan error) {
	if err := os.Remove(path); err != nil {
		errchan <- err
		return
	}
	return
}
