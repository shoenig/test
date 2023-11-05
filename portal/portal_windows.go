// Copyright (c) The Test Authors
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package portal

import (
	"net"
)

func setSocketOpt(l *net.TCPListener) error {
	// windows does not support modifying the socket; good luck!
	return nil
}
