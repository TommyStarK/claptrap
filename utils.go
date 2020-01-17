package main

import (
	"os"
	"regexp"
	"strings"
)

var parseYAMLRegexHandler = regexp.MustCompile(`(?m)\[(.+)\]=\[(.+)\]`)

func convertSignalToInt(sig os.Signal) (rc int) {
	rc = 1

	if sig == nil {
		return
	}

	switch sig.String() {
	case os.Interrupt.String():
		rc = 2
	case os.Kill.String():
		rc = 9
	case "terminated":
		rc = 15
	default:
	}

	return
}

func extractKeyValueFromStringFormatedForClaptrapYAMLConfigFile(target string) (key, value string) {
	matches := parseYAMLRegexHandler.FindStringSubmatch(target)

	if len(matches) != 3 {
		return
	}

	matches = matches[1:]
	key = strings.TrimSpace(matches[0])
	value = strings.TrimSpace(matches[1])
	return
}
