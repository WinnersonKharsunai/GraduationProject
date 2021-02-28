package queue

type Queues struct {
	Topic map[string][]Message
	DLQ   map[string][]Message
}

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

type RetrieveMessageRequest struct {
	TopicID string
}

type RetrieveMessageResponse struct {
	Message Message
}
