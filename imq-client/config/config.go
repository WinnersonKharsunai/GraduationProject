package config

// Settings contain app environment configuration settngs
type Settings struct {
	ImqClientHost       string `env:"IMQ_CLIENT_HOST" envDefault:"localhost"`
	ImqClientPort       int    `env:"IMQ_CLIENT_PORT" envDefault:"80"`
	ImqClientDailerHost string `env:"IMQ_CLIENT_DAILER_PORT" envDefault:"localhost"`
	ImqClientDailerPort int    `env:"IMQ_CLIENT_DAILER_PORT" envDefault:"5000"`
}
