package network

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type ICMPExecutor struct{}

func NewPinger() *ICMPExecutor {
	return &ICMPExecutor{}
}

func (p *ICMPExecutor) Ping(ip string, timeout time.Duration) (int64, error) {
	proto := "ip4:icmp"
	if runtime.GOOS != "windows" {
		proto = "udp4"
	}

	c, err := icmp.ListenPacket(proto, "0.0.0.0")
	if err != nil {
		return 0, fmt.Errorf("socket bind error: %w", err)
	}
	defer c.Close()

	dst, err := net.ResolveIPAddr("ip4", ip)
	if err != nil {
		return 0, fmt.Errorf("resolve error: %w", err)
	}

	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("LAG-MON"),
		},
	}
	b, _ := m.Marshal(nil)

	start := time.Now()
	if _, err := c.WriteTo(b, &net.UDPAddr{IP: dst.IP}); err != nil {
		return 0, err
	}

	if err := c.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return 0, err
	}

	rb := make([]byte, 1500)
	n, _, err := c.ReadFrom(rb)
	if err != nil {
		return 0, err
	}
	duration := time.Since(start)

	rm, err := icmp.ParseMessage(1, rb[:n])
	if err != nil || rm.Type != ipv4.ICMPTypeEchoReply {
		return 0, fmt.Errorf("invalid reply")
	}

	return duration.Microseconds(), nil
}
