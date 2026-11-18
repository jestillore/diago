package media

import (
	"net"
	"testing"
)

// ListenUDP opens both of a session's sockets, so a caller can set them up
// before they are used: bind them to one network interface, say.
func TestListenUDPOpensTheSessionsSockets(t *testing.T) {
	old := ListenUDP
	t.Cleanup(func() { ListenUDP = old })

	var opened []int
	ListenUDP = func(network string, laddr *net.UDPAddr) (*net.UDPConn, error) {
		c, err := net.ListenUDP(network, laddr)
		if err == nil {
			opened = append(opened, c.LocalAddr().(*net.UDPAddr).Port)
		}
		return c, err
	}

	s, err := NewMediaSession(net.IPv4(127, 0, 0, 1), 0)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	if len(opened) != 2 || opened[0] != s.Laddr.Port || opened[1] != s.Laddr.Port+1 {
		t.Fatalf("opened %v, want the RTP port %d and the RTCP port after it", opened, s.Laddr.Port)
	}
}

// Without a ListenUDP of its own, a session listens as it always has.
func TestListenUDPDefault(t *testing.T) {
	s, err := NewMediaSession(net.IPv4(127, 0, 0, 1), 0)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if s.Laddr.Port == 0 {
		t.Fatal("no port")
	}
}
