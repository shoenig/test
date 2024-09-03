// Copyright (c) The Test Authors
// SPDX-License-Identifier: MPL-2.0

package util_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/shoenig/test/util"
)

var t = new(testing.T)

func ExampleTempFile() {
	path := util.TempFile(t,
		util.String("hello!"),
		util.Mode(0o640),
	)

	b, _ := os.ReadFile(path)
	fmt.Println(string(b))
	// Output: hello!
}
