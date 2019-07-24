// Package logging helps to log information to Stderr
package logging

import (
	"fmt"
	"os"
)

// Debug is a flag which is set to true when debug mode is active
var Debug bool

// Println is the same as fmt.Prinln when debug mode is on
// Outputs to Stderr
func Println(args ...interface{}) {
	if !Debug {
		return
	}
	fmt.Fprintln(os.Stderr, args...)
}

// Print is the same as fmt.Print when debug mode is on
// Outputs to Stderr
func Print(args ...interface{}) {
	if !Debug {
		return
	}
	fmt.Fprint(os.Stderr, args...)
}

// Fatal is the same as log.Fatal with empty prefix string
func Fatal(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}

// Error outputs to Stderr
func Error(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}
