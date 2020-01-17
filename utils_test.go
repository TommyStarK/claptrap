package main

import (
	"os"
	"syscall"
	"testing"
)

func TestParseHTTPHeaders(t *testing.T) {
	key, value := extractKeyValueFromStringFormatedForClaptrapYAMLConfigFile("[content-type]=[application/json]")
	if key != "content-type" {
		t.Log("key should be equal to 'content-type'")
		t.Fail()
	}

	if value != "application/json" {
		t.Log("value should be equal to 'application/json'")
		t.Fail()
	}

	key, value = extractKeyValueFromStringFormatedForClaptrapYAMLConfigFile("[keep-alive]=[timeout=5, max=1000]")
	if key != "keep-alive" {
		t.Log("key should be equal to 'keep-alive'")
		t.Fail()
	}

	if value != "timeout=5, max=1000" {
		t.Log("value should be equal to 'timeout=5, max=1000'")
		t.Fail()
	}

	key, value = extractKeyValueFromStringFormatedForClaptrapYAMLConfigFile("[       foo       ]=[      bar   ]")
	if key != "foo" {
		t.Log("key should be equal to 'foo'")
		t.Fail()
	}

	if value != "bar" {
		t.Log("value should be equal to 'bar'")
		t.Fail()
	}

	key, value = extractKeyValueFromStringFormatedForClaptrapYAMLConfigFile("[invalid]=[]")
	if key != "" {
		t.Log("key should be an empty string")
		t.Fail()
	}

	if value != "" {
		t.Log("value should be an empty string")
		t.Fail()
	}
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
