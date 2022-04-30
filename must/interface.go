// Package must provides a modern generic testing assertions library.
package must

import (
	"fmt"
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

func fail(t T, msg string, args ...any) {
	c := assertions.Caller()
	s := c + fmt.Sprintf(msg, args...)
	t.Fatalf("\n" + strings.TrimSpace(s) + "\n")
}

func invoke(t T, result string) {
	if !passing(result) {
		fail(t, result)
	}
}
