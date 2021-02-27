package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/WinnersonKharsunai/GraduationProject/server/cmd/routes"
	"github.com/WinnersonKharsunai/GraduationProject/server/cmd/services/publisher"
	"github.com/WinnersonKharsunai/GraduationProject/server/cmd/services/subscriber"
	"github.com/WinnersonKharsunai/GraduationProject/server/config"
	"github.com/WinnersonKharsunai/GraduationProject/server/internal/domain"
	"github.com/WinnersonKharsunai/GraduationProject/server/internal/queue"
	"github.com/WinnersonKharsunai/GraduationProject/server/internal/storage"
	"github.com/WinnersonKharsunai/GraduationProject/server/pkg/server"
	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfgs.DbUserName, cfgs.DbPassword, cfgs.DbHost, cfgs.DbPort, cfgs.DbName)
	db, err := storage.NewMysqlDB(dsn, log)
	if err != nil {
		log.Fatal(err)
	}

	qq := queue.NewQueue(db)

	if err := qq.Init(); err != nil {
		log.Fatalln(err)
	}

	tSvc := domain.NewTopic(log, db, qq)

	pb := publisher.NewPublisher(log, tSvc)

	sc := subscriber.NewSubscriber(log, tSvc)

	handler := routes.NewHandler(pb, sc)

	addr := fmt.Sprintf("%s:%d", cfgs.ImqServerHost, cfgs.ImqServerPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	s := server.NewServer(log, lis, cfgs.PublisherCount, cfgs.SubscriberCount, handler)

	log.Infof("main: imq-server running on port: %v", addr)
	go func() {
		s.Serve()
	}()

	// graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer close(shutdown)

	select {
	case <-shutdown:
		log.Infof("main: shutting down imq-server")
		grace := time.Duration(time.Second * time.Duration(cfgs.ShutdownGrace))
		ctx, cancel := context.WithTimeout(context.Background(), grace)
		defer cancel()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			if err := qq.Shutdown(ctx); err != nil {
				log.Warnf("main: graceful shutdown failed: %v", err)
			}
			wg.Done()
		}()

		go func() {
			if err := s.Shutdown(ctx); err != nil {
				log.Warnf("main: graceful shutdown failed: %v", err)
			}
			wg.Done()
		}()
		wg.Wait()

		log.Infof("main: imq-server stopped: %v", addr)
	}

}
