package domain

type Topic struct {
	name string
	msg  string
}

type Message struct {
	MessageID string
	Data      string
	CretedAt  string
	ExpiresAt string
}
