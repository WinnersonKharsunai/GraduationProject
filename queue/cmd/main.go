package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WinnersonKharsunai/GraduationProject/queue/cmd/services"
	"github.com/WinnersonKharsunai/GraduationProject/queue/config"
	queueapi "github.com/WinnersonKharsunai/GraduationProject/queue/pkg/queue-api/go"
	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// initialize logger
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	// load configuration from environment variables
	cfgs := config.Settings{}
	if err := env.Parse(&cfgs); err != nil {
		log.Fatalf("main: failed to get configs: %v", err)
	}

	addr := fmt.Sprintf("%s:%d", cfgs.QueueHost, cfgs.QueuePort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("failed to listen:", err)
	}

	svc := services.NewQueueService(log)
	s := grpc.NewServer()
	queueapi.RegisterQueueServiceServer(s, svc)
	reflection.Register(s)

	serverErrors := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer close(shutdown)

	log.Infof("main: queue-api running on port: %v", addr)
	go func() {
		serverErrors <- s.Serve(lis)
	}()

	select {
	case <-serverErrors:
		log.Fatal("main: fatal error:", err)
	case <-shutdown:
		log.Infof("main: shutting down imq-server")
		grace := time.Duration(time.Second * time.Duration(cfgs.ShutdownGrace))
		_, cancel := context.WithTimeout(context.Background(), grace)
		defer cancel()

		log.Infof("main: imq-server stopped: %v", addr)
	}
}
