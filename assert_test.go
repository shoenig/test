package test

import (
	"fmt"
	"strings"
	"testing"
)

type testScript struct {
	label   string
	content string
}

func (ts *testScript) Label() string {
	return ts.label
}

func (ts *testScript) Content() string {
	return ts.content
}

type internalTest struct {
	t       *testing.T
	trigger bool
	helper  bool
	exp     string
	capture string
}

func (it *internalTest) PS(s string) PostScript {
	return &testScript{
		label:   "label: " + s,
		content: "content: " + s,
	}
}

func (it *internalTest) Helper() {
	it.helper = true
}

func (it *internalTest) Errorf(s string, args ...any) {
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

func (it *internalTest) post() {
	if !strings.Contains(it.capture, "PostScript |") {
		it.t.Fatal("expected post-script output")
	}
}

func newCase(t *testing.T, msg string) *internalTest {
	return &internalTest{
		t:       t,
		trigger: false,
		exp:     msg,
	}
}

func newCapture(t *testing.T) *internalTest {
	return &internalTest{
		t: t,
	}
}
