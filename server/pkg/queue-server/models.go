package queueserver

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

type Queues struct {
	Topic map[string][]Message
	DLQ   map[string][]Message
}

// type Message struct {
// 	MessageID int
// 	Data      string
// 	CretedAt  string
// 	ExpiresAt string
// }

type SendMessageRequest struct {
	TopicName string
	Message   Message
}

type DeleteMessageReqest struct {
	TopicName string
	Message   Message
}

type RetrieveMessageRequest struct {
	TopicName string
}

type RetrieveMessageResponse struct {
	Message Message
}
