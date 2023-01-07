// Package test provides a modern generic testing assertions library.
package test

import (
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

func fail(t T, msg string, scripts ...PostScript) {
	c := assertions.Caller()
	s := c + msg + "\n" + run(scripts...)
	t.Errorf("\n" + strings.TrimSpace(s) + "\n")
}

func invoke(t T, result string, settings ...Setting) {
	result = strings.TrimSpace(result)
	if !passing(result) {
		fail(t, result, scripts(settings...)...)
	}
}
