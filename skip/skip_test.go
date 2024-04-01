// Copyright (c) The Test Authors
// SPDX-License-Identifier: MPL-2.0

//go:build linux || darwin

package skip

import (
	"errors"
	"testing"
)

func TestSkip_OperatingSystem(t *testing.T) {
	OperatingSystem(t, "darwin", "linux", "windows")
	t.Fatal("expected to skip test")
}

func TestSkip_NotOperatingSystem(t *testing.T) {
	NotOperatingSystem(t, "windows")
	t.Fatal("expected to skip test")
}

func TestSkip_UserRoot(t *testing.T) {
	t.Skip("requires root")
	UserRoot(t)
	t.Fatal("expected to skip test")
}

func TestSkip_NotUserRoot(t *testing.T) {
	NotUserRoot(t)
	t.Fatal("expected to skip test")
}

func TestSkip_Architecture(t *testing.T) {
	Architecture(t, "arm64", "amd64")
	t.Fatal("expected to skip test")
}

func TestSkip_NotArchitecture(t *testing.T) {
	NotArchitecture(t, "itanium", "mips")
	t.Fatal("expected to skip test")
}

func TestSkip_DockerUnavailable(t *testing.T) {
	t.Skip("skip docker test") // gha runner

	DockerUnavailable(t)
	t.Fatal("expected to skip test")
}

func TestSkip_PodmanUnavailable(t *testing.T) {
	t.Skip("skip podman test") // gha runner

	PodmanUnavailable(t)
	t.Fatal("expected to skip test")
}

func TestSkip_CommandUnavailable(t *testing.T) {
	CommandUnavailable(t, "doesnotexist")
	t.Fatal("expected to skip test")
}

func TestSkip_MinimumCores(t *testing.T) {
	MinimumCores(t, 200)
	t.Fatal("expected to skip test")
}

func TestSkip_MaximumCores(t *testing.T) {
	MaximumCores(t, 2)
	t.Fatal("expected to skip test")
}

func TestSkip_CgroupsVersion(t *testing.T) {
	CgroupsVersion(t, 1)
	t.Fatal("expected to skip test")
}

func TestSkip_EnvironmentVariableSet(t *testing.T) {
	t.Setenv("EXAMPLE", "value")

	EnvironmentVariableSet(t, "EXAMPLE")
	t.Fatal("expected to skip test")
}

func TestSkip_EnvironmentVariableNotSet(t *testing.T) {
	EnvironmentVariableNotSet(t, "DOESNOTEXIST")
	t.Fatal("expected to skip test")
}

func TestSkip_EnvironmentVariableMatches(t *testing.T) {
	t.Setenv("EXAMPLE", "foo")

	EnvironmentVariableMatches(t, "EXAMPLE", "bar", "foo", "baz")
	t.Fatal("expected to skip test")
}

func TestSkip_EnvironmentVariableNotMatches(t *testing.T) {
	t.Setenv("EXAMPLE", "other")

	EnvironmentVariableNotMatches(t, "EXAMPLE", "bar", "foo", "baz")
	t.Fatal("expected to skip test")
}

func TestSkip_Error(t *testing.T) {
	err := errors.New("oops")
	Error(t, err)
	t.Fatal("expected to skip test")
}
