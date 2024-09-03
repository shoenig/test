// Copyright (c) The Test Authors
// SPDX-License-Identifier: MPL-2.0

package util_test

import (
	"bytes"
	"errors"
	"fmt"
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
	t.Helper()
	return t.t.TempDir()
}

func (t *helperTracker) Helper() {
	t.t.Helper()
	t.helperCalled = true
}

func (t *helperTracker) Errorf(s string, args ...any) {
	t.Helper()
	t.t.Errorf(s, args)
}

func (t *helperTracker) Fatalf(s string, args ...any) {
	t.Helper()
	t.t.Fatalf(s, args...)
}

func (t *helperTracker) Cleanup(f func()) {
	t.Helper()
	t.t.Cleanup(f)
}

func trackFailure(t util.T) *failureTracker {
	return &failureTracker{t: t}
}

type failureTracker struct {
	failed bool
	log    bytes.Buffer
	t      util.T
}

func (t *failureTracker) TempDir() string {
	t.Helper()
	return t.t.TempDir()
}

func (t *failureTracker) Helper() {
	t.t.Helper()
}

func (t *failureTracker) Errorf(s string, args ...any) {
	t.Helper()
	t.failed = true
	fmt.Fprintf(&t.log, s+"\n", args...)
}

func (t *failureTracker) Fatalf(s string, args ...any) {
	t.Helper()
	t.failed = true
	fmt.Fprintf(&t.log, s+"\n", args...)
}

func (t *failureTracker) Cleanup(f func()) {
	t.Helper()
	t.t.Cleanup(f)
}

func (t *failureTracker) AssertFailedWith(msg string) {
	t.Helper()
	if !t.failed {
		t.t.Fatalf("expected test to fail with message %q", msg)
	}
	strlog := t.log.String()
	if !strings.Contains(strlog, msg) {
		t.t.Fatalf("expected test to fail with message %q\ngot message %q", msg, strlog)
	}
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

	t.Run("using multiple options", func(t *testing.T) {
		var expectedMode fs.FileMode = 0444
		expectedData := "important data"
		prefix := "harvey-"
		pattern := prefix + "*"
		path := util.TempFile(t,
			util.Mode(expectedMode),
			util.Pattern(pattern),
			util.String(expectedData))

		t.Run("Mode sets a custom file mode", func(t *testing.T) {
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
			if !strings.Contains(path, prefix) {
				t.Fatalf("filename does not match pattern\nexpected to contain %s\ngot %s", prefix, path)
			}
		})
		t.Run("sets string data", func(t *testing.T) {
			actualData, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read temp file: %v", err)
			}
			if expectedData != string(actualData) {
				t.Fatalf("temp file contains wrong data\nexpected %q\ngot %q", expectedData, string(actualData))
			}
		})
	})

	t.Run("Bytes sets binary data", func(t *testing.T) {
		expectedData := []byte("important data")
		path := util.TempFile(t, util.Bytes(expectedData))
		actualData, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read temp file: %v", err)
		}
		if !bytes.Equal(expectedData, actualData) {
			t.Fatalf("temp file contains wrong data\nexpected %q\ngot %q", string(expectedData), actualData)
		}
	})

	t.Run("cleans up file (and nothing else) in custom dir", func(t *testing.T) {
		dir := t.TempDir()
		existing, err := os.CreateTemp(dir, "")
		if err != nil {
			t.Fatalf("failed to create temporary file: %v", err)
		}
		existingPath := existing.Name()
		existing.Close()
		var newPath string

		t.Run("uses custom directory", func(t *testing.T) {
			newPath = util.TempFile(t, util.Dir(dir))
			entries, err := os.ReadDir(dir)
			if err != nil {
				t.Fatalf("failed to read directory: %v", err)
			}
			var found bool
			for _, entry := range entries {
				if entry.Name() == filepath.Base(newPath) {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("did not find temporary file in %s", dir)
			}
		})

		if newPath == "" {
			t.Fatal("expected non-empty path")
		}
		_, err = os.Stat(newPath)
		if !errors.Is(err, fs.ErrNotExist) {
			if err == nil {
				t.Errorf("expected temp file not to exist: %s", newPath)
			} else {
				t.Errorf("unexpected error: %v", err)
			}
		}
		_, err = os.Stat(existingPath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				t.Error("expected pre-existing file not to be deleted")
			} else {
				t.Errorf("unexpected error statting pre-existing file: %v", err)
			}
		}
	})

	t.Run("fails if specified dir doesn't exist", func(t *testing.T) {
		fakeDir := filepath.Join(t.TempDir(), "fake")
		tracker := trackFailure(t)
		path := util.TempFile(tracker, util.Dir(fakeDir))
		if path != "" {
			t.Errorf("expected empty path\ngot %q", path)
		}
		tracker.AssertFailedWith("TempFile: directory does not exist")
	})
}
