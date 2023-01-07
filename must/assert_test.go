package must

import (
	"fmt"
	"strings"
)

func (it *internalTest) Fatalf(s string, args ...any) {
	if !it.trigger {
		it.trigger = true
	}
	msg := strings.TrimSpace(fmt.Sprintf(s, args...))
	it.capture = msg
	fmt.Println(msg)
}
