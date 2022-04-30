package must

import (
	"fmt"
	"strings"
	"testing"
)

type internalTest struct {
	t       *testing.T
	trigger bool
	helper  bool
	exp     string
	capture string
}

func (it *internalTest) Helper() {
	it.helper = true
}

func (it *internalTest) Fatalf(s string, args ...any) {
	if !it.trigger {
		it.trigger = true
	}
	msg := strings.TrimSpace(fmt.Sprintf(s, args...))
	it.capture = msg
	fmt.Println(msg)
}

func (it *internalTest) assert() {
	if !it.helper {
		it.t.Fatal("should be marked as helper")
	}
	if !it.trigger {
		it.t.Fatalf("condition expected to trigger; did not")
	}

	if !strings.Contains(it.capture, it.exp) {
		it.t.Fatalf("expected message %q in output, got %q", it.exp, it.capture)
	}
}

func newCase(t *testing.T, msg string) *internalTest {
	return &internalTest{
		t:       t,
		trigger: false,
		exp:     msg,
	}
}
