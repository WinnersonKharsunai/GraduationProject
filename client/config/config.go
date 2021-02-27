package config

import (
	"os"
	"strconv"
)

// Settings contain app environment configuration settngs
type Settings struct {
	ImqClientHost       string `env:"IMQ_CLIENT_HOST" envDefault:"localhost"`
	ImqClientPort       int    `env:"IMQ_CLIENT_PORT" envDefault:"80"`
	ShutdownGrace       int    `env:"SHUTDOWN_GRACE" envDefault:"10"`
	ImqClientDailerHost string `env:"IMQ_CLIENT_DAILER_PORT" envDefault:"localhost"`
	ImqClientDailerPort int    `env:"IMQ_CLIENT_DAILER_PORT" envDefault:"80"`
}

// GetDailerEndpoint ...
func GetDailerEndpoint() (string, int) {
	endpoint := os.Args
	host := endpoint[0]
	port, _ := strconv.Atoi(endpoint[1])
	return host, port
}
