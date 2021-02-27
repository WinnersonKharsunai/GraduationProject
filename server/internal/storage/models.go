package storage

type Message struct {
	MessageID int
	Data      string
	CretedAt  string
	ExpiresAt string
}

type Queue struct {
	Topic map[string][]Message
	DLQ   map[string][]Message
}
