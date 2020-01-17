package main

import (
	"testing"
)

func TestParseHTTPHeaders(t *testing.T) {
	key, value := mapHTTPHeaders("[content-type]=[application/json]")
	if key != "content-type" {
		t.Log("header key should be equal to 'content-type'")
		t.Fail()
	}

	if value != "application/json" {
		t.Log("header value should be equal to 'application/json'")
		t.Fail()
	}

	key, value = mapHTTPHeaders("[keep-alive]=[timeout=5, max=1000]")
	if key != "keep-alive" {
		t.Log("header key should be equal to 'keep-alive'")
		t.Fail()
	}

	if value != "timeout=5, max=1000" {
		t.Log("header value should be equal to 'timeout=5, max=1000'")
		t.Fail()
	}

	key, value = mapHTTPHeaders("[       foo       ]=[      bar   ]")
	if key != "foo" {
		t.Log("header key should be equal to 'foo'")
		t.Fail()
	}

	if value != "bar" {
		t.Log("header value should be equal to 'bar'")
		t.Fail()
	}
}
