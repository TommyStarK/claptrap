package main

import (
	"fmt"
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
		testDataPath = "./testdata"
	}

	fmt.Println(onCI)
	fmt.Println(gopath)
	fmt.Println(testDataPath)
}

func TestClaptrapInstanciationShouldFail(t *testing.T) {
	if _, err := newClaptrap("invalid", nil); err == nil {
		t.Log("provided invalid path, should have failed to instanciate claptrap")
		t.Fail()
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

	go c.trap()

	triggerWrite := make(chan chan struct{})
	go func() {
		writeDone := <-triggerWrite
		writeBigFile(testDataPath+"/bigfile", testContent, c.errors)
		writeDone <- struct{}{}
		return
	}()

	triggerUpdate := make(chan chan struct{})
	go func() {
		updateDone := <-triggerUpdate
		writeFile(testDataPath+"/bigfile", testContent, c.errors)
		updateDone <- struct{}{}
		return
	}()

	triggerRename := make(chan chan struct{})
	go func() {
		renameDone := <-triggerRename
		renameFile(testDataPath+"/bigfile", testDataPath+"/bigf", c.errors)
		renameDone <- struct{}{}
		return
	}()

	triggerRemove := make(chan chan struct{})
	go func() {
		removeDone := <-triggerRemove
		removeFile(testDataPath+"/bigf", c.errors)
		removeDone <- struct{}{}
		return
	}()

	witness := make(chan struct{})
	triggerWrite <- witness
	<-witness
	close(triggerWrite)
	processResult("CREATE", "testdata/bigfile", ch, t)

	triggerUpdate <- witness
	<-witness
	close(triggerUpdate)
	processResult("UPDATE", "testdata/bigfile", ch, t)

	triggerRename <- witness
	<-witness
	close(triggerRename)
	processResult("RENAME", "testdata/bigfile", ch, t)

	triggerRemove <- witness
	<-witness
	close(triggerRemove)
	processResult("REMOVE", "testdata/bigf", ch, t)

	c.sigchan <- os.Signal(syscall.SIGTERM)

	close(ch)
	close(witness)
}

func TestConvertSignalToInt(t *testing.T) {
	var (
		sigint  = os.Signal(syscall.SIGINT)
		sigkill = os.Signal(syscall.SIGKILL)
		sigterm = os.Signal(syscall.SIGTERM)
	)

	if convertSignalToInt(sigint) != 2 {
		t.Logf("return code should be equal to 2")
		t.Fail()
	}
	if convertSignalToInt(sigkill) != 9 {
		t.Logf("return code should be equal to 9")
		t.Fail()
	}
	if convertSignalToInt(sigterm) != 15 {
		t.Logf("return code should be equal to 15")
		t.Fail()
	}
	if convertSignalToInt(nil) != 1 {
		t.Logf("return code should be equal to 1")
		t.Fail()
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func processResult(expectedAction, expectedTarget string, ch chan [3]string, t *testing.T) {
	result := <-ch
	action, target, timestamp := result[0], result[1], result[2]

	if len(timestamp) == 0 {
		t.Log("timestamp should not be empty")
		t.Fail()
		return
	}
	if action != expectedAction || target != expectedTarget {
		t.Logf("event caught should be '%s' and target '%s' but got: [%s|%s] ", expectedAction, expectedTarget, action, target)
		t.Fail()
		return
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

	return
}
