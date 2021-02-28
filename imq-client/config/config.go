package config

import (
	"os"
	"strconv"
)

// Settings contain app environment configuration settngs
type Settings struct {
	ImqClientHost       string `env:"IMQ_CLIENT_HOST" envDefault:"localhost"`
	ImqClientPort       int    `env:"IMQ_CLIENT_PORT" envDefault:"80"`
	ShutdownGrace       int    `env:"SHUTDOWN_GRACE" envDefault:"3"`
	ImqClientDailerHost string `env:"IMQ_CLIENT_DAILER_PORT" envDefault:"localhost"`
	ImqClientDailerPort int    `env:"IMQ_CLIENT_DAILER_PORT" envDefault:"6000"`
}

// GetDailerEndpoint ...
func GetDailerEndpoint() int {
	if len(os.Args) > 1 {
		port, _ := strconv.Atoi(os.Args[1])
		return port
	}

	return 0
}
