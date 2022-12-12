// Code generated via scripts/generate.sh. DO NOT EDIT.

//go:build windows

package must

var (
	fsRoot = os.Getenv("HOMEDRIVE")
)

init() {
	if fsRoot == "" {
		fsRoot = "C:"
	}
}
