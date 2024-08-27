// Copyright (c) The Test Authors
// SPDX-License-Identifier: MPL-2.0

package util_test

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shoenig/test/util"
)

func trackHelper(t util.T) *helperTracker {
	return &helperTracker{t: t}
}

type helperTracker struct {
	helperCalled bool
	t            util.T
}

func (t *helperTracker) TempDir() string {
	t.t.Helper()
	return t.t.TempDir()
}

func (t *helperTracker) Helper() {
	t.t.Helper()
	t.helperCalled = true
}

func (t *helperTracker) Errorf(s string, args ...any) {
	t.t.Helper()
	t.t.Errorf(s, args)
}

func (t *helperTracker) Fatalf(s string, args ...any) {
	t.t.Helper()
	t.t.Fatalf(s, args...)
}

func (t *helperTracker) Cleanup(f func()) {
	t.t.Helper()
	t.t.Cleanup(f)
}

func TestTempFile(t *testing.T) {
	t.Run("creates a read/write temp file by default", func(t *testing.T) {
		th := trackHelper(t)
		path := util.TempFile(th)
		if !th.helperCalled {
			t.Errorf("expected TempFile to call Helper")
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("failed to stat temp file: %v", err)
		}
		mode := info.Mode()
		if mode&0400 == 0 || mode&0200 == 0 {
			t.Fatalf("expected at least u+rw permission, got %03o", mode)
		}
	})
	t.Run("sets a custom file mode", func(t *testing.T) {
		var expectedMode fs.FileMode = 0444
		path := util.TempFile(t, util.Mode(expectedMode))
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("failed to stat temp file: %v", err)
		}
		actualMode := info.Mode()
		if expectedMode != actualMode {
			t.Fatalf("file has wrong mode\nexpected %03o\ngot %03o", expectedMode, actualMode)
		}
	})
	t.Run("sets a name pattern", func(t *testing.T) {
		prefix := "harvey-"
		pattern := prefix + "*"
		path := util.TempFile(t, util.Pattern(pattern))
		if !strings.Contains(path, prefix) {
			t.Fatalf("filename does not match pattern\nexpected to contain %s\ngot %s", prefix, path)
		}
	})
	t.Run("sets string data", func(t *testing.T) {
		expectedData := "important data"
		path := util.TempFile(t, util.StringData(expectedData))
		actualData, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read temp file: %v", err)
		}
		if expectedData != string(actualData) {
			t.Fatalf("temp file contains wrong data\nexpected %q\ngot %q", expectedData, string(actualData))
		}
	})
	t.Run("sets binary data", func(t *testing.T) {
		expectedData := []byte("important data")
		path := util.TempFile(t, util.ByteData(expectedData))
		actualData, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read temp file: %v", err)
		}
		if !bytes.Equal(expectedData, actualData) {
			t.Fatalf("temp file contains wrong data\nexpected %q\ngot %q", string(expectedData), actualData)
		}
	})

	t.Run("file is deleted after test", func(t *testing.T) {
		dirpath := t.TempDir()
		var path string

		t.Run("uses custom path", func(t *testing.T) {
			path = util.TempFile(t, util.Path(dirpath))
			entries, err := os.ReadDir(dirpath)
			if err != nil {
				t.Fatalf("failed to read directory: %v", err)
			}
			if entries[0].Name() != filepath.Base(path) {
				t.Fatalf("did not find temporary file in %s", dirpath)
			}
		})

		if path == "" {
			t.Fatal("expected non-empty path")
		}
		_, err := os.Stat(path)
		if err == nil {
			t.Fatalf("expected temp file not to exist: %s", path)
		}
	})
}
