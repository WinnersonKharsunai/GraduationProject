package config

// Settings contain app environment configuration settngs
type Settings struct {
	ImqServerHost   string `env:"IMQ_SERVER_HOST" envDefault:"localhost"`
	ImqServerPort   int    `env:"IMQ_SERVER_PORT" envDefault:"80"`
	PublisherCount  int    `env:"PUBLISHER_COUNT" envDefault:"2"`
	SubscriberCount int    `env:"SUBSCRIBER_COUNT" envDefault:"2"`
	ShutdownGrace   int    `env:"SHUTDOWN_GRACE" envDefault:"10"`

	ImqQueueHost string `env:"IMQ_QUEUE_HOST" envDefault:""`
	ImqQueuePort int    `env:"IMQ_QUEUE_PORT" envDefault:""`

	DbUserName string `env:"DB_USERNAME" envDefault:"root"`
	DbPassword string `env:"DB_PASSWORD" envDefault:"password"`
	DbHost     string `env:"DB_HOST" envDefault:"localhost"`
	DbPort     int    `env:"DB_PORT" envDefault:"3308"`
	DbName     string `env:"DB_NAME" envDefault:"messagingQueueDev"`
}
