// Copyright (c) The Test Authors
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"errors"
	"fmt"
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
	dir         *string
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

// Dir specifies a directory path to contain the temporary file.
// If dir is the empty string, the file will be created in the
// default directory for temporary files, as returned by os.TempDir.
// A temporary file created in a custom directory will still be deleted
// after the test runs, though the directory may not.
func Dir(dir string) TempFileSetting {
	return func(s *TempFileSettings) {
		s.dir = &dir
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
	path, err := tempFile(t.Helper, t.TempDir, settings...)
	t.Cleanup(func() {
		err := os.Remove(path)
		if err != nil {
			t.Fatalf("failed to clean up temp file: %s", path)
		}
	})
	if err != nil {
		t.Fatalf("%v", err)
	}
	return path
}

// tempFile returns errors instead of relying upon T to stop execution, for ease
// of testing TempFile.
func tempFile(helper func(), tempDir func() string, settings ...TempFileSetting) (path string, err error) {
	helper()
	var allSettings TempFileSettings
	for _, setting := range settings {
		setting(&allSettings)
	}
	if allSettings.mode == nil {
		allSettings.mode = new(fs.FileMode)
		*allSettings.mode = 0600
	}
	if allSettings.dir == nil {
		allSettings.dir = new(string)
		*allSettings.dir = tempDir()
	}

	wrap := func(err error) error {
		return fmt.Errorf("TempFile: %w", err)
	}
	file, err := os.CreateTemp(*allSettings.dir, allSettings.namePattern)
	if errors.Is(err, fs.ErrNotExist) {
		return "", fmt.Errorf("TempFile: directory does not exist")
	}
	if err != nil {
		return "", wrap(err)
	}
	path = file.Name()
	_, err = file.Write(allSettings.data)
	if err != nil {
		file.Close()
		return path, wrap(err)
	}
	err = file.Close()
	if err != nil {
		return path, wrap(err)
	}
	err = os.Chmod(path, *allSettings.mode)
	if err != nil {
		return path, wrap(err)
	}
	return file.Name(), nil
}
