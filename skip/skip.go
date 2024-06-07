// Copyright (c) The Test Authors
// SPDX-License-Identifier: MPL-2.0

// Package skip provides helper functions for skipping test cases.
package skip

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// T is the minimal set of functions to be implemented by any testing
// framework compatible with the skip package.
type T interface {
	Skipf(string, ...any)
	Fatalf(string, ...any)
}

// OperatingSystem will skip the test if the Go runtime detects the operating
// system matches one of the given names.
func OperatingSystem(t T, names ...string) {
	os := runtime.GOOS
	for _, name := range names {
		if os == strings.ToLower(name) {
			t.Skipf("operating system excluded from tests %q", os)
		}
	}
}

// NotOperatingSystem will skip the test if the Go runtime detects the operating
// system does not match one of the given names.
func NotOperatingSystem(t T, names ...string) {
	os := runtime.GOOS
	for _, name := range names {
		if os == strings.ToLower(name) {
			return
		}
	}
	t.Skipf("operating excluded from tests %q", os)
}

// UserRoot will skip the test if the test is being run as the root user.
//
// Uses the effective UID value to determine user.
func UserRoot(t T) {
	euid := os.Geteuid()
	if euid == 0 {
		t.Skipf("test must not run as root")
	}
}

// NotUserRoot will skip the test if the test is not being run as the root user.
//
// Uses the effective UID value to determine user.
func NotUserRoot(t T) {
	euid := os.Geteuid()
	if euid != 0 {
		t.Skipf("test must run as root")
	}
}

// Architecture will skip the test if the Go runtime detects the system
// architecture matches one of the given names.
func Architecture(t T, names ...string) {
	arch := runtime.GOARCH
	for _, name := range names {
		if arch == strings.ToLower(name) {
			t.Skipf("arch excluded from tests %q", arch)
		}
	}
}

// NotArchitecture will skip the test if the Go runtime the system architecture
// does not match one of the given names.
func NotArchitecture(t T, names ...string) {
	arch := runtime.GOARCH
	for _, name := range names {
		if arch == strings.ToLower(name) {
			return
		}
	}
	t.Skipf("arch excluded from tests %q", arch)
}

func cmdAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return !errors.Is(err, exec.ErrNotFound)
}

// CommandUnavailable will skip the test if the given command cannot be found on
// the system PATH.
func CommandUnavailable(t T, command string) {
	if !cmdAvailable(command) {
		t.Skipf("command %q not detected on system", command)
	}
}

// DockerUnavailable will skip the test if the docker command cannot be found on
// the system PATH.
func DockerUnavailable(t T) {
	if !cmdAvailable("docker") {
		t.Skipf("docker not detected on system")
	}
}

// PodmanUnavailable will skip the test if the podman command cannot be found on
// the system PATH.
func PodmanUnavailable(t T) {
	if !cmdAvailable("podman") {
		t.Skipf("podman not detected on system")
	}
}

// MinimumCores will skip the test if the system does not meet the minimum
// number of CPU cores.
func MinimumCores(t T, num int) {
	cpus := runtime.NumCPU()
	if cpus < num {
		t.Skipf("system does not meet minimum cpu cores")
	}
}

// MaximumCores will skip the test if the number of cores on the system
// exceeds the given maximum.
func MaximumCores(t T, num int) {
	cpus := runtime.NumCPU()
	if cpus > num {
		t.Skipf("system exceeds maximum cpu cores")
	}
}

// CgroupsVersion will skip the test if the system does not match the given
// cgroups version.
func CgroupsVersion(t T, version int) {
	if runtime.GOOS != "linux" {
		t.Skipf("cgroups requires linux")
	}

	mType := mountType(t, "/sys/fs/cgroup")

	switch mType {
	case "tmpfs":
		// this is a cgroups v1 system
		if version == 2 {
			t.Skipf("system does not match cgroups version 2")
		}
	case "cgroup2":
		// this is a cgroups v2 system
		if version == 1 {
			t.Skipf("system does not match cgroups version 1")
		}
	default:
		t.Fatalf("unknown cgroups mount type %q", mType)
	}
}

func mountType(t T, path string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "df", "-T", path)
	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unable to run df command: %v", err)
	}

	// need the first token of the second line
	output := string(b)
	tokenRe := regexp.MustCompile(`on\s+([\w]+)\s+`)
	results := tokenRe.FindStringSubmatch(output)
	if len(results) != 2 {
		t.Fatalf("no mount type for path")
	}
	return results[1]
}

// EnvironmentVariableSet will skip the test if the given environment variable
// is set to any value.
func EnvironmentVariableSet(t T, name string) {
	if name == "" {
		t.Fatalf("environment variable name must be set")
	}

	_, exists := os.LookupEnv(name)
	if exists {
		t.Skipf("environment variable %q is set", name)
	}
}

// EnvironmentVariableNotSet will skip the test if the given environment
// variable is not set in the system environment.
func EnvironmentVariableNotSet(t T, name string) {
	if name == "" {
		t.Fatalf("environment variable name must be set")
	}
	_, exists := os.LookupEnv(name)
	if !exists {
		t.Skipf("environment variable %q is not set", name)
	}
}

// EnvironmentVariableMatches will skip the test if the system environment
// matches one of the given environment variable values.
func EnvironmentVariableMatches(t T, name string, values ...string) {
	if len(values) == 0 {
		t.Fatalf("no possible environment variable values given")
	}

	// if not set in the system, then it must not match
	actual, exists := os.LookupEnv(name)
	if !exists {
		return
	}

	for _, value := range values {
		if value == actual {
			t.Skipf("environment variable %q matches %q", name, value)
		}
	}
}

// EnvironmentVariableNotMatches will skip the test if the system environment
// does not match one of the given environment variable values.
func EnvironmentVariableNotMatches(t T, name string, values ...string) {
	if len(values) == 0 {
		t.Fatalf("no possible environment variable values given")
	}

	actual, exists := os.LookupEnv(name)
	if !exists {
		t.Skipf("environment variable %q not set", name)
	}

	for _, value := range values {
		if actual == value {
			return
		}
	}

	t.Skipf("environment variable %q does not match values (is %q)", name, actual)
}

// Error will skip the test if err is not nil.
func Error(t T, err error) {
	if err != nil {
		t.Skipf("skipping test due to non-nil error: %v", err)
	}
}
