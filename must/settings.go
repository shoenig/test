// Code generated via scripts/generate.sh. DO NOT EDIT.

package must

// Settings are used to manage a collection of Setting values used to modify
// the behavior of a test case assertion.
type Settings struct {
	postScripts []PostScript
}

// A Setting changes the behavior of a test case assertion.
type Setting func(s *Settings)
