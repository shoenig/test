//go:build windows

package test

import (
	"os"
)

var (
	fsRoot = os.Getenv("HOMEDRIVE")
)

func init() {
	if fsRoot == "" {
		fsRoot = "C:"
	}
}
