package main

import "fmt"

// Handle has the expected signature by claptrap.
// Arguments will always be given in that order. Your logic
// implemented in the following function will be executed in a
// separate goroutine for each event caught by claptrap.
func Handle(action string, target string, timestamp string) {
	fmt.Printf("event [%s] detected at [%s] target: %s\n", action, timestamp, target)
	return
}
