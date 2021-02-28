package server

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/WinnersonKharsunai/GraduationProject/server/cmd/routes"
	"github.com/sirupsen/logrus"
)

// Server is the concrete implementation for the Multiclient hanlder
type Server struct {
	log             *logrus.Logger
	lis             net.Listener
	router          routes.Router
	publisherCount  int
	subscriberCount int
	clientType      map[string]string
	newConnection   chan newConnection
	processCh       chan net.Conn
	processWg       sync.WaitGroup
	shutdownCh      chan struct{}
}

// NewServer is the factory function for the Server type
func NewServer(log *logrus.Logger, lis net.Listener, pubCount, subCount int, router routes.Router) *Server {
	s := &Server{
		log:             log,
		lis:             lis,
		publisherCount:  pubCount,
		subscriberCount: subCount,
		router:          router,
		clientType:      make(map[string]string),
		shutdownCh:      make(chan struct{}),
		newConnection:   make(chan newConnection),
		processCh:       make(chan net.Conn),
	}

	s.processWg.Add(pubCount)
	for i := 1; i <= pubCount; i++ {
		s.clientType[fmt.Sprintf("%v", 4999+i)] = pub
		go s.processWorker(i)
	}

	s.processWg.Add(subCount)
	for i := 1; i <= pubCount; i++ {
		s.clientType[fmt.Sprintf("%v", 5999+i)] = sub
		go s.processWorker(i)
	}

	return s
}

// Serve starts the server
func (s *Server) Serve() error {
	shutdown := false
	for !shutdown {
		select {
		case <-s.shutdownCh:
			shutdown = true
		default:
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
			defer cancel()

			go func() {
				s.accept(ctx)
			}()

			select {
			case <-ctx.Done():
			case ch := <-s.newConnection:
				if ch.err != nil {
					s.log.Errorf("failed to accept client: %v", ch.err)
					continue
				}
				s.processCh <- ch.conn
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
	}
}

// Shutdown gracefully shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		close(s.shutdownCh)

		close(s.processCh)
		s.processWg.Wait()

		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
