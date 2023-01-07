package test

import "testing"

func TestPostScript_Label(t *testing.T) {
	var s PostScript = &script{label: "\nhello "}
	if s.Label() != "hello" {
		t.FailNow()
	}
}

func TestPostScript_Content(t *testing.T) {
	var s PostScript = &script{content: "\nhello "}
	if s.Content() != "\thello" {
		t.FailNow()
	}
}

func TestPostScript_Sprintf(t *testing.T) {
	ps := Sprintf("foo %s %d", "baz", 1)
	result := run(scripts(ps)...)
	exp := "↪ PostScript | annotation ↷\n\tfoo baz 1\n"
	if result != exp {
		t.Fatalf("exp %s, got %s", exp, result)
	}
}

func TestPostScript_KV(t *testing.T) {
	ps := Values("one", 1, "foo", "bar")
	result := run(scripts(ps)...)
	exp := "↪ PostScript | mapping ↷\n\t\"one\" => 1\n\t\"foo\" => \"bar\"\n"
	if result != exp {
		t.Fatalf("exp %s, got %s", exp, result)
	}
}

func TestPostScript_Func(t *testing.T) {
	ps := Func(func() string {
		return "hello"
	})
	result := run(scripts(ps)...)
	exp := "↪ PostScript | function ↷\n\thello\n"
	if result != exp {
		t.Fatalf("exp %s, got %s", exp, result)
	}
}
