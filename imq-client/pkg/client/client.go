package client

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/WinnersonKharsunai/GraduationProject/client/pkg/protocol"
)

// Client is the concrete implementation for the client
type Client struct {
	Addr       string
	con        net.Conn
	dialerHost string
	dialerPort int
}

// Service is the interface for the client
type Service interface {
	Dial() error
	GetID() int
	GetAddress() string
	SendRequest(ctx context.Context, request *protocol.Request) ([]byte, error)
}

// NewClient is the factory functipn for the client
func NewClient(addr, dialerHost string, dialerPort int) Service {
	return &Client{
		Addr:       addr,
		dialerHost: dialerHost,
		dialerPort: dialerPort,
	}
}

// Dial helps dail for a new connection
func (c *Client) Dial() error {
	dialer := &net.Dialer{
		LocalAddr: &net.TCPAddr{
			IP:   net.ParseIP(c.dialerHost),
			Port: c.dialerPort,
		},
	}

	con, err := dialer.Dial("tcp", c.Addr)
	if err != nil {
		return err
	}
	c.con = con

	if err = testConnection(con); err != nil {
		return err
	}

	return nil
}

// GetAddress return cleint address
func (c *Client) GetAddress() string {
	return c.Addr
}

// GetID return client Id
func (c *Client) GetID() int {
	return c.dialerPort
}

// SendRequest send the request to server
func (c *Client) SendRequest(ctx context.Context, request *protocol.Request) ([]byte, error) {

	if err := writeToConnection(c.con, request); err != nil {
		return []byte{}, err
	}

	raw, err := readFromConnection(c.con)
	if err != nil {
		return []byte{}, err
	}

	response, err := unmarshalResponse(raw)
	if err != nil {
		return []byte{}, err
	}

	if response.Error != "" {
		return []byte{}, errors.New(response.Error)
	}

	return response.Body, nil
}

func testConnection(c net.Conn) error {
	data, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		return err
	}

	response, err := unmarshalResponse(data)
	if err != nil {
		return err
	}

	if response.Error != "" {
		return errors.New(response.Error)
	}

	return nil
}

func readFromConnection(c net.Conn) (string, error) {
	data, err := bufio.NewReader(c).ReadString('\n')
	return data, err
}

func writeToConnection(c net.Conn, request *protocol.Request) error {
	b, err := json.Marshal(request)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(c, "%v\n", string(b))
	return err
}

func unmarshalResponse(raw string) (*protocol.Response, error) {
	response := protocol.Response{}

	err := json.Unmarshal([]byte(raw), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
