package server

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/WinnersonKharsunai/GraduationProject/server/pkg/protocol"
)

func (s *Server) processWorker(id int) {
	for con := range s.processCh {
		s.log.Infof("processWorker with Id %v: get a new connection with: %v", id, con.RemoteAddr().String())

		response := &protocol.Response{}
		connectionClosed := false

		if s.clientType[getClientPort(con)] == "" {
			err := errors.New("unauthorised")
			s.log.Errorf("processWorker with Id %v: failed to recognised client: %v", id, err)
			response.Error = err.Error()
			connectionClosed = true
		}

		response.Body = []byte(statusConnected)

		if err := writeToConnection(con, response); err != nil {
			s.log.Errorf("processWorker with Id %v: %v:%v", id, failedTowriteResponse, err)
		}

		for !connectionClosed {
			data, err := readFromConnection(con)
			if err != nil {
				s.log.Errorf("processWorker with Id %v: failed to read client request: %v", id, err)
				connectionClosed = true
				continue
			}

			response = s.router.RequestRouter(context.Background(), data)

			if err := writeToConnection(con, response); err != nil {
				s.log.Errorf("processWorker with Id %v: %v:%v", id, failedTowriteResponse, err)
			}
		}

		s.log.Infof("processWorker with Id %v: closing connection with: %v", id, con.RemoteAddr().String())
		con.Close()
	}
	s.processWg.Done()
}

func getClientPort(c net.Conn) string {
	remoteAddr := strings.SplitAfter(c.RemoteAddr().String(), ":")
	return remoteAddr[1]
}

func readFromConnection(c net.Conn) (string, error) {
	data, err := bufio.NewReader(c).ReadString('\n')
	return data, err
}

func writeToConnection(c net.Conn, response *protocol.Response) error {
	b, err := json.Marshal(response)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(c, "%v\n", string(b))
	return err
}
