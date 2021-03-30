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
	log := initializeLogger()

	cfgs, err := loadConfigurations()
	if err != nil {
		log.Fatalf("main: failed to get configs: %v", err)
	}

	db, err := connectToDatabase(cfgs)
	if err != nil {
		log.Fatalf("main: failed to connect to database: %v", err)
	}

	queueSvc, err := startQueue(db, log)
	if err != nil {
		log.Fatalf("main: failed to start queue service: %v", err)
	}

	handler := initializeServiceHandler(log, db, queueSvc)

	serverr, addr, err := startImqServer(log, cfgs, handler)
	if err != nil {
		log.Fatalf("main: failed to start listener: %v", err)
	}

	gracefulShutdown(log, addr, serverr, queueSvc, cfgs.ShutdownGrace)
}

func initializeLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	return log
}

func loadConfigurations() (config.Settings, error) {
	cfgs := config.Settings{}
	if err := env.Parse(&cfgs); err != nil {
		return config.Settings{}, err
	}
	return cfgs, nil
}

func connectToDatabase(cfgs config.Settings) (storage.DatabaseIF, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfgs.DbUserName, cfgs.DbPassword, cfgs.DbHost, cfgs.DbPort, cfgs.DbName)
	db, err := storage.NewMysqlDB(dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func startQueue(db storage.DatabaseIF, log *logrus.Logger) (queue.ImqQueueIF, error) {
	queueSvc, err := queue.NewQueue(log, db)
	if err != nil {
		return nil, err
	}
	return queueSvc, nil
}

func initializeServiceHandler(log *logrus.Logger, db storage.DatabaseIF, qSvc queue.ImqQueueIF) routes.Router {
	topicSvc := domain.NewTopic(log, db, qSvc)
	publisherSvc := publisher.NewPublisher(log, topicSvc)
	subscriberSvc := subscriber.NewSubscriber(log, topicSvc)
	return routes.NewHandler(publisherSvc, subscriberSvc)
}

func startImqServer(log *logrus.Logger, cfgs config.Settings, handler routes.Router) (*server.Server, string, error) {
	addr := fmt.Sprintf("%s:%d", cfgs.ImqServerHost, cfgs.ImqServerPort)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, addr, err
	}

	s := server.NewServer(log, lis, cfgs.PublisherCount, cfgs.SubscriberCount, handler)
	log.Infof("main: imq-server running on port: %v", addr)

	go func() {
		s.Serve()
	}()

	return s, addr, nil
}

func gracefulShutdown(log *logrus.Logger, addr string, servr *server.Server, queueSvc queue.ImqQueueIF, shutdownGrace int) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer close(shutdown)

	select {
	case <-shutdown:
		log.Infof("main: shutting down imq-server")
		grace := time.Duration(time.Second * time.Duration(shutdownGrace))
		ctx, cancel := context.WithTimeout(context.Background(), grace)
		defer cancel()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			if err := queueSvc.BackUpQueue(ctx); err != nil {
				log.Warnf("main: graceful shutdown failed: %v", err)
			}
			wg.Done()
		}()

		go func() {
			if err := servr.Shutdown(ctx); err != nil {
				log.Warnf("main: graceful shutdown failed: %v", err)
			}
			wg.Done()
		}()
		wg.Wait()

		log.Infof("main: imq-server stopped: %v", addr)
	}
}
