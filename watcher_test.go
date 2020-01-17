package main

import (
	"os"
	"testing"
	"time"

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

func TestWatcherShortEvent(t *testing.T) {
	evch := make(chan *event)
	errch := make(chan error)

	testWatcher, err := newWatcher("./testdata", evch, errch)
	if err != nil {
		t.Fatal(err)
	}

	go testWatcher.watch()

	go func() {
		time.Sleep(1 * time.Second)
		writeFile("./testdata/foo", testContent, errch)
		return
	}()

	go func() {
		time.Sleep(3 * time.Second)
		renameFile("./testdata/foo", "./testdata/bar", errch)
		return
	}()

	go func() {
		time.Sleep(5 * time.Second)
		removeFile("./testdata/bar", errch)
		return
	}()

	results := make([]*event, 0)

Loop:
	for {
		select {
		case <-time.After(10 * time.Second):
			t.Fatal("test timed out")
		case e, ok := <-evch:
			if !ok || e == nil {
				t.Fatal()
			}

			results = append(results, e)
			if len(results) == 3 {
				break Loop
			}

		case err, ok := <-errch:
			if !ok || err == nil {
				t.Fatal()
			}

			t.Fatal(err)
		}
	}

	if len(results) != 3 {
		t.Fatal("expecting 3 events: new file, rename file and remove file")
	}

	if timestamp, ok := results[1].trace[fsnotify.Rename]; !ok || len(timestamp) == 0 {
		t.Log("second event caught should be: RENAME")
		t.Fail()
	}

	if timestamp, ok := results[2].trace[fsnotify.Remove]; !ok || len(timestamp) == 0 {
		t.Log("last event caught should be: REMOVE")
		t.Fail()
	}

	if err := testWatcher.stop(); err != nil {
		t.Log(err)
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
