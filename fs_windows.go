//go:build windows

package test

var (
	fsRoot = os.Getenv("HOMEDRIVE")
)

init() {
	if fsRoot == "" {
		fsRoot = "C:"
	}
}
