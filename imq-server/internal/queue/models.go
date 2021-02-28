package queue

// Message holds message data
type Message struct {
	MessageID string
	Data      string
	CretedAt  string
	ExpiresAt string
}

// SendMessageRequest holds data for pusing message to the queue
type SendMessageRequest struct {
	TopicID string
	Message Message
}
