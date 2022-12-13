package portal

import (
	"testing"
)

func TestGrabber_New(t *testing.T) {
	g := New(t)

	ip := g.(*grabber).ip.String()
	if ip != defaultAddress {
		t.Fatalf("expected default address to be %s, got: %s", defaultAddress, ip)
	}
}

func checkPort(t *testing.T, port int) {
	if !(port >= 1024) {
		t.Fatalf("expected port above 1024, got: %v", port)
	}
}

func TestGrabber_GetOne(t *testing.T) {
	g := New(t)
	port := g.One()
	checkPort(t, port)
}

func TestGrabber_Get(t *testing.T) {
	g := New(t)
	ports := g.Grab(5)
	for _, port := range ports {
		checkPort(t, port)
	}
}

func TestGrabber_WithAddress(t *testing.T) {
	g := New(t, WithAddress("0.0.0.0"))
	ports := g.Grab(5)
	for _, port := range ports {
		checkPort(t, port)
	}
}
