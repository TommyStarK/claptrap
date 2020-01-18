package main

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func init() {

}

func TestClaptrapInstanciationShouldFail(t *testing.T) {
	if _, err := newClaptrap("invalid", nil); err == nil {
		t.Log("provided invalid path, should have failed to instanciate claptrap")
		t.Fail()
	}
}

func writeBigFile(path, content string, errchan chan error) {
	var f *os.File
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		errchan <- err
		return
	}

	for i := 0; i < 1000; i++ {
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

func TestWriteBigFile(t *testing.T) {
	c, err := newClaptrap("./testdata", nil)
	if err != nil {
		t.Fatal(err)
	}

	c.testMode = true
	go func() {
		time.Sleep(1 * time.Second)
		writeBigFile("./testdata/bigfile", testContent, c.errors)
		return
	}()

	go func() {
		time.Sleep(10*time.Second - 2*time.Second)
		removeFile("./testdata/bigfile", c.errors)
		return
	}()

	go func() {
		time.Sleep(10 * time.Second)
		c.sigchan <- os.Signal(syscall.SIGTERM)
		return
	}()

	c.trap()
}

func TestSendSIGTERM(t *testing.T) {
	c, err := newClaptrap("./testdata", nil)
	if err != nil {
		t.Fatal(err)
	}

	c.testMode = true
	go func() {
		time.Sleep(1 * time.Second)
		c.sigchan <- os.Signal(syscall.SIGTERM)
		return
	}()

	c.trap()
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
