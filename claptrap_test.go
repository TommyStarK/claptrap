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

	if convertSignalToInt(sigint) != 2+128 {
		t.Logf("return code should be equal to 2+128")
		t.Fail()
	}

	if convertSignalToInt(sigkill) != 9+128 {
		t.Logf("return code should be equal to 9+128")
		t.Fail()
	}

	if convertSignalToInt(sigterm) != 15+128 {
		t.Logf("return code should be equal to 15+128")
		t.Fail()
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
