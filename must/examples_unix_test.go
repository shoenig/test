// Code generated via scripts/generate.sh. DO NOT EDIT.

// Copyright (c) The Test Authors
// SPDX-License-Identifier: MPL-2.0

//go:build unix

package must

import (
	"io/fs"
	"os"
)

func ExampleDirExists() {
	DirExists(t, "/tmp")
	// Output:
}

func ExampleDirNotExists() {
	DirNotExists(t, "/does/not/exist")
	// Output:
}

func ExampleFileContains() {
	_ = os.WriteFile("/tmp/example", []byte("foo bar baz"), fs.FileMode(0600))
	FileContains(t, "/tmp/example", "bar")
	// Output:
}

func ExampleFileExists() {
	_ = os.WriteFile("/tmp/example", []byte{}, fs.FileMode(0600))
	FileExists(t, "/tmp/example")
	// Output:
}

func ExampleFileMode() {
	_ = os.WriteFile("/tmp/example_fm", []byte{}, fs.FileMode(0600))
	FileMode(t, "/tmp/example_fm", fs.FileMode(0600))
	// Output:
}

func ExampleFileNotExists() {
	FileNotExists(t, "/tmp/not_existing_file")
	// Output:
}
