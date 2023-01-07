package test

import (
	"github.com/google/go-cmp/cmp"
)

// Settings are used to manage a collection of Setting values used to modify
// the behavior of a test case assertion.
type Settings struct {
	postScripts []PostScript
	cmpOptions  []cmp.Option
}

// A Setting changes the behavior of a test case assertion.
type Setting func(s *Settings)

// Cmp enables configuring cmp.Option values for test case assertions that make
// use of the cmp package.
func Cmp(options ...cmp.Option) Setting {
	return func(s *Settings) {
		s.cmpOptions = append(s.cmpOptions, options...)
	}
}

func options(settings ...Setting) []cmp.Option {
	s := new(Settings)
	for _, setting := range settings {
		setting(s)
	}
	return s.cmpOptions
}

func scripts(settings ...Setting) []PostScript {
	s := new(Settings)
	for _, setting := range settings {
		setting(s)
	}
	return s.postScripts
}
