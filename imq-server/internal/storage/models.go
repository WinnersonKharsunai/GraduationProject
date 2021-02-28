package storage

type Message struct {
	MessageID string
	Data      string
	CretedAt  string
	ExpiresAt string
}

type Queue struct {
	Topic map[string][]Message
}

type StoreQueue struct {
	QueuID    string
	TopicID   string
	MessageID string
}
