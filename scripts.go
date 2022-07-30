package test

import (
	"fmt"
	"strings"
)

func run(posts ...PostScript) string {
	s := new(strings.Builder)
	for _, post := range posts {
		s.WriteString("↪ PostScript | ")
		s.WriteString(post.Label())
		s.WriteString(" ↷\n")
		s.WriteString(post.Content())
		s.WriteString("\n")
	}
	return s.String()
}

// A PostScript is used to annotate a test failure with additional information.
//
// Can be useful in large e2e style test cases, where adding additional context
// beyond an assertion helps in debugging.
type PostScript interface {
	// Label should categorize what is in Content.
	Label() string

	// Content contains extra contextual information for debugging a test failure.
	Content() string
}

type script struct {
	label   string
	content string
}

func (s *script) Label() string {
	return strings.TrimSpace(s.label)
}
func (s *script) Content() string {
	return "\t" + strings.TrimSpace(s.content)
}

// Sprintf appends a Sprintf-string as an annotation to the output of a test case failure.
func Sprintf(msg string, args ...any) PostScript {
	return &script{
		label:   "annotation",
		content: fmt.Sprintf(msg, args...),
	}
}

// Sprint appends a Sprint-string as an annotation to the output of a test case failure.
func Sprint(args ...any) PostScript {
	return &script{
		label:   "annotation",
		content: strings.TrimSpace(fmt.Sprintln(args...)),
	}
}

// Values adds formatted key-value mappings as an annotation to the output of a test case failure.
func Values(values ...any) PostScript {
	b := new(strings.Builder)
	n := len(values)
	for i := 0; i < n-1; i += 2 {
		s := fmt.Sprintf("\t%#v => %#v\n", values[i], values[i+1])
		b.WriteString(s)
	}
	if n%2 != 0 {
		s := fmt.Sprintf("\t%v => <MISSING ARG>", values[n-1])
		b.WriteString(s)
	}
	return &script{
		label:   "mapping",
		content: b.String(),
	}
}

// Func adds the string produced by f as an annotation to the output of a test case failure.
func Func(f func() string) PostScript {
	return &script{
		label:   "function",
		content: f(),
	}
}
