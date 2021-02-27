package client

import (
	"bufio"
	"errors"
	"net"
)

// Client ...
type Client struct {
	addr       string
	dialerHost string
	dialerPort int
}

// NewClient ...
func NewClient(addr, dialerHost string, dialerPort int) *Client {
	return &Client{
		addr:       addr,
		dialerHost: dialerHost,
		dialerPort: dialerPort,
	}
}

// Dial ...
func (c *Client) Dial() error {
	dialer := &net.Dialer{
		LocalAddr: &net.TCPAddr{
			IP:   net.ParseIP(c.dialerHost),
			Port: c.dialerPort,
		},
	}

	con, err := dialer.Dial("tcp", c.addr)
	if err != nil {
		return err
	}

	return testConnection(con)
}

func testConnection(c net.Conn) error {
	response, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		return err
	}

	if response != "connected" {
		return errors.New(response)
	}
	return nil
}
