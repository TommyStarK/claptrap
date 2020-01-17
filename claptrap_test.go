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
	var cfg = &config{
		path: "invalid",
	}

	if _, err := newClaptrap(cfg); err == nil {
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
	var cfg = &config{
		path: "./testdata",
	}

	c, err := newClaptrap(cfg)
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
	var cfg = &config{
		path: "./testdata",
	}

	c, err := newClaptrap(cfg)
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

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
