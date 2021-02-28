package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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
	// initialize logger
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	// load configuration from environment variables
	cfgs := config.Settings{}
	if err := env.Parse(&cfgs); err != nil {
		log.Fatalf("failed to get configs: %v", err)
	}

	port := config.GetDailerEndpoint()

	if port > 0 {
		cfgs.ImqClientDailerPort = port
	}

	addr := fmt.Sprintf("%s:%d", cfgs.ImqClientHost, cfgs.ImqClientPort)
	c := client.NewClient(addr, cfgs.ImqClientDailerHost, cfgs.ImqClientDailerPort)

	log.Infof("Dailing for server connection")
	if err := c.Dial(); err != nil {
		log.Fatalln(err)
	}
	log.Infof("connected")

	mf := messagefactory.NewMessageFactory()

	p := publisher.NewPublisher(c, mf)
	s := subscriber.NewSubscriber(c, mf)

	cs := console.NewConsole(cfgs.ImqClientDailerPort, c, p, s)

	err := cs.Start(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

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
		wg.Add(1)

		go func() {
			if err := cs.Shutdown(ctx); err != nil {
				log.Warnf("main: graceful shutdown failed: %v", err)
			}
			wg.Done()
		}()
		wg.Wait()

		log.Infof("main: imq-server stopped: %v", addr)
	}
}
