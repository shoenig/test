// Package must provides a modern generic testing assertions library.
package must

import (
	"strings"

	"github.com/shoenig/test/internal/assertions"
)

// T is the minimal set of functions to be implemented by any testing framework
// compatible with the must package.
type T interface {
	Helper()
	Fatalf(string, ...any)
}

func passing(result string) bool {
	return result == ""
}

func fail(t T, msg string, scripts ...PostScript) {
	c := assertions.Caller()
	s := c + msg + run(scripts...)
	t.Fatalf("\n" + strings.TrimSpace(s) + "\n")
}

func invoke(t T, result string, scripts ...PostScript) {
	if !passing(result) {
		fail(t, result, scripts...)
	}
}
