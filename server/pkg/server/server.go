package server

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	publisher "github.com/WinnersonKharsunai/GraduationProject/server/cmd/services/publisher"
	subscriber "github.com/WinnersonKharsunai/GraduationProject/server/cmd/services/subscriber"
	"github.com/WinnersonKharsunai/GraduationProject/server/pkg/protocol"
	"github.com/sirupsen/logrus"
)

type newConnection struct {
	conn net.Conn
	addr string
	err  error
}

// Server ...
type Server struct {
	log             *logrus.Logger
	lis             net.Listener
	pSvc            publisher.PublisherIF
	sSvc            subscriber.SubscriberIF
	publisherCount  int
	subscriberCount int
	protocol        protocol.Protocol
	newConnection   chan newConnection
	publisherCh     chan net.Conn
	subscriberCh    chan net.Conn
	publisherWg     sync.WaitGroup
	subscriberWg    sync.WaitGroup
	shutdownCh      chan struct{}
}

// NewServer ...
func NewServer(log *logrus.Logger, lis net.Listener, pubCount, subCount int, pSvc publisher.PublisherIF, sSvc subscriber.SubscriberIF) *Server {
	s := &Server{
		log:             log,
		lis:             lis,
		publisherCount:  pubCount,
		subscriberCount: subCount,
		pSvc:            pSvc,
		sSvc:            sSvc,
		shutdownCh:      make(chan struct{}),
		newConnection:   make(chan newConnection),
		publisherCh:     make(chan net.Conn, pubCount),
		subscriberCh:    make(chan net.Conn, subCount),
	}

	s.publisherWg.Add(pubCount)
	for i := 1; i <= pubCount; i++ {
		go s.publisherWorker(i)
	}

	s.subscriberWg.Add(subCount)
	for i := 1; i <= pubCount; i++ {
		go s.subscriberWorker(i)
	}

	return s
}

// Serve ...
func (s *Server) Serve() error {
	shutdown := false
	for !shutdown {
		select {
		case <-s.shutdownCh:
			shutdown = true
		default:
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
			defer cancel()

			go s.accept(ctx)

			select {
			case <-ctx.Done():
			case ch := <-s.newConnection:
				if ch.err != nil {
					s.log.Errorf("failed to accept client: %v", ch.err)
					continue
				}
				s.log.Infof("got a new connection: %v", ch.addr)
				port := strings.SplitAfter(ch.addr, ":")
				switch port[1] {
				case "5000":
					s.publisherCh <- ch.conn
				case "6000":
					s.subscriberCh <- ch.conn
				default:
					s.log.Infof("given port not expected: %v", ch.addr)
				}
			}
		}
	}
	return nil
}

func (s *Server) accept(ctx context.Context) {
	con, err := s.lis.Accept()
	s.newConnection <- newConnection{
		conn: con,
		err:  err,
		addr: con.RemoteAddr().String(),
	}
}

func (s *Server) publisherWorker(i int) {
	s.log.Infof("publisherWorker with Id %v: started", i)

	for con := range s.publisherCh {
		s.log.Infof("publisherWorker with Id %v: get a new connection with: %v", i, con.RemoteAddr().String())

		done := false
		for !done {
			data, err := readFromConn(con)
			if err != nil {
				s.log.Errorf("publisherWorker with Id %v: failed to read client data :%v", i, err)
				done = true
				continue
			}

			var response protocol.Response
			resp, err := s.processRequest(context.Background(), data)
			if err != nil {
				response.Error = err.Error()
			}

			response.Body = resp

			if err := writeResponse(con, response); err != nil {
				s.log.Info("publisherWorker with Id %v: failed to write response: %v", i, err)
			}

			s.log.Infof("publisherWorker with Id %v: response: %v", i, response)
		}
		con.Close()
	}
	s.publisherWg.Done()
}

func (s *Server) subscriberWorker(i int) {
	s.log.Infof("subscriberWorker with Id %v: started", i)

	for con := range s.subscriberCh {
		s.log.Infof("subscriberWorker with Id %v: get a new connection with: %v", i, con.RemoteAddr().String())

		done := false
		for !done {
			data, err := readFromConn(con)
			if err != nil {
				s.log.Errorf("publisherWorker with Id %v: failed to read client data :%v", i, err)
				done = true
				continue
			}

			var response protocol.Response
			resp, err := s.processRequest(context.Background(), data)
			if err != nil {
				response.Error = err.Error()
			}

			response.Body = resp

			if err := writeResponse(con, response); err != nil {
				s.log.Info("publisherWorker with Id %v: failed to write response: %v", i, err)
			}

			s.log.Infof("publisherWorker with Id %v: response: %v", i, response)
		}
		con.Close()
	}
	s.subscriberWg.Done()
}

func readFromConn(con net.Conn) (string, error) {
	data, err := bufio.NewReader(con).ReadString('\n')
	return data, err
}

func writeResponse(con net.Conn, response protocol.Response) error {
	_, err := fmt.Fprintf(con, "%v\n", response)
	if err != nil {
		return err
	}
	return nil
}

func parseRequest(rawRequest string) (*protocol.Request, error) {
	r := strings.Split(rawRequest, "|")
	if len(r) != 2 {
		return nil, errors.New("Bad request")
	}

	var hdr protocol.Header
	if err := json.Unmarshal([]byte(r[0]), &hdr); err != nil {
		return nil, err
	}

	req := protocol.Request{
		Header: hdr,
		Body:   r[1],
	}

	return &req, nil
}

func (s *Server) processRequest(ctx context.Context, rawRequest string) (string, error) {
	request, err := parseRequest(rawRequest)
	if err != nil {
		return "", err
	}

	if err := s.protocol.ValidateRequestHeader(ctx, request.Header); err != nil {
		return "", err
	}

	switch request.Header.Method {
	case "ShowTopics":
		_, err := s.pSvc.ShowTopics(ctx)
		if err != nil {
			return "", err
		}
		return "", nil
	default:
		return "", errors.New("Unknown method type")
	}
}

// Shutdown gracefully shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		close(s.shutdownCh)

		close(s.publisherCh)
		s.publisherWg.Wait()

		close(s.subscriberCh)
		s.subscriberWg.Wait()

		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
