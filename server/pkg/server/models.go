package server

import "net"

const (
	pub                   = "publisher"
	sub                   = "subscriber"
	failedTowriteResponse = "failed to write response to client"
)

type newConnection struct {
	conn net.Conn
	err  error
}
