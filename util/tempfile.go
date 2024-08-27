// Copyright (c) The Test Authors
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"io/fs"
	"os"
)

type T interface {
	TempDir() string
	Helper()
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	Cleanup(func())
}

type TempFileSettings struct {
	data        []byte
	mode        *fs.FileMode
	namePattern string
	path        *string
}

type TempFileSetting func(s *TempFileSettings)

// Pattern sets the filename to pattern with a random string appended.
// If pattern contains a '*', the last '*' will be replaced by the
// random string.
func Pattern(pattern string) TempFileSetting {
	return func(s *TempFileSettings) {
		s.namePattern = pattern
	}
}

// Mode sets the temporary file's mode.
func Mode(mode fs.FileMode) TempFileSetting {
	return func(s *TempFileSettings) {
		s.mode = &mode
	}
}

// StringData writes data to the temporary file.
func StringData(data string) TempFileSetting {
	return func(s *TempFileSettings) {
		s.data = []byte(data)
	}
}

// ByteData writes data to the temporary file.
func ByteData(data []byte) TempFileSetting {
	return func(s *TempFileSettings) {
		s.data = data
	}
}

// Path specifies a directory path to contain the temporary file.
// A temporary file created in a custom directory will still be deleted
// after the test runs, though the directory may not.
func Path(path string) TempFileSetting {
	return func(s *TempFileSettings) {
		s.path = &path
	}
}

// TempFile creates a temporary file that is deleted after the test is
// completed. If the file cannot be deleted, the test fails with a message
// containing its path. TempFile creates a new file every time it is called.
// By default, each file thus created is in a unique directory as
// created by (*testing.T).TempDir(); this directory is also deleted
// after the test is completed.
func TempFile(t T, settings ...TempFileSetting) (path string) {
	t.Helper()
	var allSettings TempFileSettings
	for _, setting := range settings {
		setting(&allSettings)
	}
	if allSettings.mode == nil {
		allSettings.mode = new(fs.FileMode)
		*allSettings.mode = 0600
	}
	if allSettings.path == nil {
		allSettings.path = new(string)
		*allSettings.path = t.TempDir()
	}

	var err error
	crash := func(t T) 
		t.Helper()
		t.Fatalf("%s: %v", "TempFile", err)
	}
	file, err := os.CreateTemp(*allSettings.path, allSettings.namePattern)
	if err != nil {
		crash()
	}
	path = file.Name()
	t.Cleanup(func() {
		err := os.Remove(path)
		if err != nil {
			t.Fatalf("failed to clean up temp file: %s", path)
		}
	})
	_, err = file.Write(allSettings.data)
	if err != nil {
		file.Close()
		crash()
	}
	err = file.Close()
	if err != nil {
		crash()
	}
	err = os.Chmod(path, *allSettings.mode)
	if err != nil {
		crash()
	}
	return file.Name()
}
