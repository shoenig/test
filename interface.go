// Package test provides a modern generic testing assertions library.
package test

import (
	"fmt"
	"strings"

	"github.com/shoenig/test/internal/assertions"
)

// T is the minimal set of functions to be implemented by any testing framework
// compatible with the test package.
type T interface {
	Helper()
	Errorf(string, ...any)
}

func passing(result string) bool {
	return result == ""
}

func fail(t T, msg string, args ...any) {
	c := assertions.Caller()
	s := c + fmt.Sprintf(msg, args...)
	t.Errorf("\n" + strings.TrimSpace(s) + "\n")
}

func invoke(t T, result string) {
	if !passing(result) {
		fail(t, result)
	}
}
