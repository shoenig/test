//go:build windows

package portal

import (
	"net"
)

func setSocketOpt(l *net.TCPListener) error {
	// windows does not support modifying the socket; good luck!
	return nil
}
