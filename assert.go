// Copyright (c) The Test Authors
// SPDX-License-Identifier: MPL-2.0

package test

// T is the minimal set of functions to be implemented by any testing framework
// compatible with the test package.
type T interface {
	Helper()
	Errorf(string, ...any)
}

func errorf(t T, msg string, args ...any) {
	t.Errorf(msg, args...)
}
