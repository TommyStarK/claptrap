package main

import (
	"regexp"
	"strings"
)

var parseHTTPHeaderRegex = regexp.MustCompile(`(?m)\[(.+)\]=\[(.+)\]`)

func mapHTTPHeaders(header string) (key, value string) {
	matches := parseHTTPHeaderRegex.FindStringSubmatch(header)

	if len(matches) != 3 {
		return
	}

	matches = matches[1:]
	key = strings.TrimSpace(matches[0])
	value = strings.TrimSpace(matches[1])
	return
}
