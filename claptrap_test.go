package main

import (
	"os"
	"runtime"
	"syscall"
	"testing"
)

var (
	onCI         = len(os.Getenv("CI")) > 0
	gopath       = os.Getenv("GOPATH")
	testDataPath = gopath + "/src/github.com/TommyStarK/claptrap/testdata"
)

func init() {
	if len(gopath) == 0 {
		panic("$GOPATH not set")
	}

	if !onCI {
		testDataPath = "testdata"
	}
}

func TestClaptrapInstanciationShouldFail(t *testing.T) {
	if _, err := newClaptrap("invalid", nil); err == nil {
		t.Fatal("provided invalid path, should have failed to instanciate claptrap")
	}
}

func TestClaptrapBehaviorOnLargeFile(t *testing.T) {
	c, err := newClaptrap(testDataPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan [3]string)
	c.testMode = true
	c.testchan = ch

	triggerWrite := make(chan chan struct{}, 1)
	triggerUpdate := make(chan chan struct{}, 1)
	triggerRename := make(chan chan struct{}, 1)
	triggerRemove := make(chan chan struct{}, 1)
	witness := make(chan struct{}, 1)

	go c.trap()

	go func(triggerWrite chan chan struct{}) {
		writeDone := <-triggerWrite
		writeBigFile(testDataPath+"/bigfile", testContent, c.errors)
		writeDone <- struct{}{}
	}(triggerWrite)

	go func(triggerUpdate chan chan struct{}) {
		updateDone := <-triggerUpdate
		writeFile(testDataPath+"/bigfile", testContent, c.errors)
		updateDone <- struct{}{}
	}(triggerUpdate)

	go func(triggerRename chan chan struct{}) {
		renameDone := <-triggerRename
		renameFile(testDataPath+"/bigfile", testDataPath+"/bigf", c.errors)
		renameDone <- struct{}{}
	}(triggerRename)

	go func(triggerRemove chan chan struct{}) {
		removeDone := <-triggerRemove
		removeFile(testDataPath+"/bigf", c.errors)
		removeDone <- struct{}{}
	}(triggerRemove)

	triggerWrite <- witness
	<-witness
	close(triggerWrite)
	processResult(t, "CREATE", testDataPath+"/bigfile", ch)

	triggerUpdate <- witness
	<-witness
	close(triggerUpdate)
	processResult(t, "UPDATE", testDataPath+"/bigfile", ch)

	triggerRename <- witness
	<-witness
	close(triggerRename)
	processResult(t, "RENAME", testDataPath+"/bigfile", ch)

	triggerRemove <- witness
	<-witness
	close(triggerRemove)
	processResult(t, "REMOVE", testDataPath+"/bigf", ch)

	c.sigchan <- os.Signal(syscall.SIGTERM)

	close(ch)
	close(witness)
}

func TestConvertSignalToInt(t *testing.T) {
	for _, test := range []struct {
		name     string
		signal   os.Signal
		expected int
	}{
		{
			name:     "nil",
			signal:   nil,
			expected: 1,
		},
		{
			name:     "not supported",
			signal:   syscall.SIGABRT,
			expected: 1,
		},
		{
			name:     "sigint",
			signal:   syscall.SIGINT,
			expected: 2,
		},
		{
			name:     "sigkill",
			signal:   syscall.SIGKILL,
			expected: 9,
		},
		{
			name:     "sigterm",
			signal:   syscall.SIGTERM,
			expected: 15,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if c := convertSignalToInt(test.signal); c != test.expected {
				t.Errorf("expected: %d, but got: %d", test.expected, c)
			}
		})
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func processResult(t *testing.T, expectedAction, expectedTarget string, ch chan [3]string) {
	result := <-ch
	action, target, timestamp := result[0], result[1], result[2]

	if len(timestamp) == 0 {
		t.Fatal("timestamp should not be empty")
	}
	if action != expectedAction || target != expectedTarget {
		t.Fatalf("event caught should be '%s' and target '%s' but got: [%s|%s] ", expectedAction, expectedTarget, action, target)
	}
}

func writeBigFile(path, content string, errchan chan error) {
	var f *os.File
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		errchan <- err
		return
	}

	for i := 0; i < runtime.NumCPU()*1000; i++ {
		if _, err := f.WriteString(content); err != nil {
			errchan <- err
			return
		}
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
