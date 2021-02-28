package server

import "net"

type newConnection struct {
	conn net.Conn
	err  error
}

const (
	pub                   = "publisher"
	sub                   = "subscriber"
	statusConnected       = "connected"
	failedTowriteResponse = "failed to write response to client"
)
