package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/WinnersonKharsunai/GraduationProject/client/cmd/services/publisher"
	"github.com/WinnersonKharsunai/GraduationProject/client/cmd/services/subscriber"
	"github.com/WinnersonKharsunai/GraduationProject/client/config"
	"github.com/WinnersonKharsunai/GraduationProject/client/console"
	messagefactory "github.com/WinnersonKharsunai/GraduationProject/client/message-factory"
	"github.com/WinnersonKharsunai/GraduationProject/client/pkg/client"
	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

func main() {
	log := initializeLogger()

	cfgs, err := loadConfigurations()
	if err != nil {
		log.Fatalf("main: failed to get configs: %v", err)
	}

	port := getDailerEndpoint()
	if port > 0 {
		cfgs.ImqClientDailerPort = port
	}

	log.Infof("main: dailing for server connection")
	clientSvc, addr, err := connectToServer(cfgs)
	if err != nil {
		log.Fatalf("main: failed to connect to server: %v", err)
	}
	log.Infof("main: connected to server")

	consoleSvc := initializeConsole(clientSvc, cfgs.ImqClientDailerPort)

	err = consoleSvc.Start(context.Background())
	if err != nil {
		log.Fatalf("main: failed to start console: %v", err)
	}

	gracefulShutdown(log, addr, consoleSvc)
}

func getDailerEndpoint() int {
	if len(os.Args) > 1 {
		port, _ := strconv.Atoi(os.Args[1])
		return port
	}
	return 0
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

func connectToServer(cfgs config.Settings) (client.Service, string, error) {
	addr := fmt.Sprintf("%s:%d", cfgs.ImqClientHost, cfgs.ImqClientPort)

	c := client.NewClient(addr, cfgs.ImqClientDailerHost, cfgs.ImqClientDailerPort)
	if err := c.Dial(); err != nil {
		return nil, addr, err
	}
	return c, addr, nil
}

func initializeConsole(clientSvc client.Service, clientID int) *console.Console {
	messagefactorySvc := messagefactory.NewMessageFactory()
	publisherSvc := publisher.NewPublisher(clientSvc, messagefactorySvc)
	subscriberSvc := subscriber.NewSubscriber(clientSvc, messagefactorySvc)
	return console.NewConsole(clientID, clientSvc, publisherSvc, subscriberSvc)
}

func gracefulShutdown(log *logrus.Logger, addr string, consoleSvc *console.Console) {
	defer close(consoleSvc.ShutdwonChan)
	select {
	case <-consoleSvc.ShutdwonChan:
		log.Infof("main: imq-server stopped: %v", addr)
	}
}
