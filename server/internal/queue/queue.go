package queue

import (
	"context"

	"github.com/WinnersonKharsunai/GraduationProject/server/internal/storage"
)

// ImqQueueIF is the inteerface for the Queue
type ImqQueueIF interface {
	Init() error
	SendMessage(ctx context.Context, message SendMessageRequest) error
	RetrieveMessage(ctx context.Context, topicName string) Message
	DeleteMessage(ctx context.Context, m DeleteMessageReqest)
	Shutdown(ctx context.Context) error
}

// Queue is the concrete implementztion for Queue
type Queue struct {
	db                        storage.DatabaseIF
	queueChan                 chan struct{}
	sendMessageChan           chan SendMessageRequest
	retrieveMessageChan       chan RetrieveMessageRequest
	retrieveMessageResponseCh chan RetrieveMessageResponse
	deleteChan                chan DeleteMessageReqest
	shutdownChan              chan struct{}
}

// NewQueue is the factory function for the Queue
func NewQueue(db storage.DatabaseIF) ImqQueueIF {
	return &Queue{
		db:                        db,
		queueChan:                 make(chan struct{}),
		sendMessageChan:           make(chan SendMessageRequest),
		retrieveMessageChan:       make(chan RetrieveMessageRequest),
		retrieveMessageResponseCh: make(chan RetrieveMessageResponse),
		deleteChan:                make(chan DeleteMessageReqest),
		shutdownChan:              make(chan struct{}),
	}
}

// Init load Queue from the database
func (q *Queue) Init() error {

	queue, err := q.loadQueues(context.Background())
	if err != nil {
		return err
	}

	go q.queueService(*queue)

	return nil
}

// SendMessage push message to the queue
func (q *Queue) SendMessage(ctx context.Context, message SendMessageRequest) error {
	q.sendMessageChan <- SendMessageRequest{
		TopicID: message.TopicID,
		Message: message.Message,
	}
	return nil
}

// RetrieveMessage pull message from the queue
func (q *Queue) RetrieveMessage(ctx context.Context, topicName string) Message {
	q.retrieveMessageChan <- RetrieveMessageRequest{
		TopicName: topicName,
	}

	msg := <-q.retrieveMessageResponseCh

	return Message{
		MessageID: msg.Message.MessageID,
		Data:      msg.Message.Data,
		CretedAt:  msg.Message.CretedAt,
		ExpiresAt: msg.Message.ExpiresAt,
	}
}

// DeleteMessage pop message from the queue
func (q *Queue) DeleteMessage(ctx context.Context, m DeleteMessageReqest) {
	q.deleteChan <- DeleteMessageReqest{
		TopicName: m.TopicName,
		Message:   m.Message,
	}
}

// Shutdown gracefully shutdown the Queue Service
func (q *Queue) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
