// Package portal (Port Allocator) provides a helper for reserving free TCP ports
// across multiple processes on the same machine. This works by asking the kernel
// for available ports in the ephemeral port range. It does so by binding to an
// address with port 0 (e.g. 127.0.0.1:0), modifying the socket to disable SO_LINGER,
// close the connection, and finally return the port that was used. This *probably*
// works well, because the kernel re-uses ports in an LRU fashion, implying the
// test code asking for the ports *should* be the only thing immediately asking
// to bind that port again.
package portal

import (
	"io"
	"net"
	"strconv"
	"sync"
)

const (
	defaultAddress = "127.0.0.1"
)

type FatalTester interface {
	Fatalf(msg string, args ...any)
}

// A Grabber is used to grab open ports.
type Grabber interface {
	// Grab n port allocations.
	Grab(n int) []int

	// One port allocation.
	One() int
}

// New creates a new Grabber with the given options.
func New(t FatalTester, opts ...Option) Grabber {
	g := &grabber{
		t:  t,
		ip: net.ParseIP(defaultAddress),
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

type grabber struct {
	t    FatalTester
	ip   net.IP
	lock sync.Mutex
}

type Option func(Grabber)

// WithAddress specifies which address on which to allocate ports.
func WithAddress(address string) Option {
	return func(g Grabber) {
		g.(*grabber).ip = net.ParseIP(address)
	}
}

func (g *grabber) Grab(n int) []int {
	g.lock.Lock()
	defer g.lock.Unlock()

	ports := make([]int, n)
	closers := make([]io.Closer, n)

	for i := 0; i < n; i++ {
		p, c := g.one()
		ports[i] = p
		closers[i] = c
	}

	for _, c := range closers {
		_ = c.Close()
	}

	return ports
}

func (g *grabber) One() int {
	g.lock.Lock()
	defer g.lock.Unlock()

	p, c := g.one()
	_ = c.Close()
	return p
}

// one will acquire one port; the caller must hold the lock and also close
// the returned listener - this minimized the chances of reallocating the same
// port
func (g *grabber) one() (int, io.Closer) {
	tcpAddr := &net.TCPAddr{IP: g.ip, Port: 0}
	l, listenErr := net.ListenTCP("tcp", tcpAddr)
	if listenErr != nil {
		g.t.Fatalf("failed to acquire port: %v", listenErr)
	}

	if setErr := setSocketOpt(l); setErr != nil {
		g.t.Fatalf("failed to modify socket: %v", setErr)
	}

	_, port, splitErr := net.SplitHostPort(l.Addr().String())
	if splitErr != nil {
		g.t.Fatalf("failed to parse address: %v", splitErr)
	}

	p, parseErr := strconv.Atoi(port)
	if parseErr != nil {
		g.t.Fatalf("failed to parse port: %v", parseErr)
	}

	return p, l
}
