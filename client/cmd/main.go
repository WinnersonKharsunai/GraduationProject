package main

import (
	"fmt"

	"github.com/WinnersonKharsunai/GraduationProject/client/config"
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

	addr := fmt.Sprintf("%s:%d", cfgs.ImqClientHost, cfgs.ImqClientPort)
	c := client.NewClient(addr, cfgs.ImqClientDailerHost, cfgs.ImqClientDailerPort)

	if err := c.Dial(); err != nil {
		log.Fatalln()
	}
}
