package main

import (
	"os"
	"syscall"
	"testing"
)

func init() {

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
