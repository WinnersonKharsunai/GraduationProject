package server

import "net"

const (
	pub                   = "publisher"
	sub                   = "subscriber"
	failedTowriteResponse = "failed to write response to client"
	statusConnected       = "connected"
)

type newConnection struct {
	conn net.Conn
	err  error
}
