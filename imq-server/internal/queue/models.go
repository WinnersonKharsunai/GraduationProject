package queue

type Message struct {
	MessageID string
	Data      string
	CretedAt  string
	ExpiresAt string
}

type SendMessageRequest struct {
	TopicID string
	Message Message
}
