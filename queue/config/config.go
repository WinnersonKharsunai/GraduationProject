package config

// Settings contains app environment settings
type Settings struct {
	QueueHost     string `env:"QUEUE_HOST" envDefault:"localhost"`
	QueuePort     int    `env:"QUEUE_PORT" envDefault:"9999"`
	ShutdownGrace int    `env:"SHUTDOWN_GRACE" envDefault:"10"`

	DbUserName string `env:"DB_USERNAME" envDefault:"root"`
	DbPassword string `env:"DB_PASSWORD" envDefault:"password"`
	DbHost     string `env:"DB_HOST" envDefault:"localhost"`
	DbPort     int    `env:"DB_PORT" envDefault:"3308"`
	DbName     string `env:"DB_NAME" envDefault:"queue"`
}
