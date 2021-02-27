package queueserver

import (
	"context"
	"database/sql"
	"fmt"
)

// QueueServer ...
type QueueServer struct {
	Dsn                       string
	Cxn                       *sql.DB
	queueChan                 chan struct{}
	sendMessageChan           chan SendMessageRequest
	retrieveMessageChan       chan RetrieveMessageRequest
	retrieveMessageResponseCh chan RetrieveMessageResponse
	deleteChan                chan DeleteMessageReqest
	shutdownChan              chan struct{}
}

// NewQueueServer ...
func NewQueueServer(dataSourseName string) *QueueServer {
	qs := QueueServer{
		Dsn:                       dataSourseName,
		queueChan:                 make(chan struct{}),
		sendMessageChan:           make(chan SendMessageRequest),
		retrieveMessageChan:       make(chan RetrieveMessageRequest),
		retrieveMessageResponseCh: make(chan RetrieveMessageResponse),
		deleteChan:                make(chan DeleteMessageReqest),
		shutdownChan:              make(chan struct{}),
	}

	return &qs
}

func (q *QueueServer) queueService(queue Queues) {

	fmt.Println(queue)

	shutdown := false
	for !shutdown {
		select {
		case <-q.shutdownChan:
			if err := q.saveQueues(context.Background(), &queue); err != nil {
				fmt.Println(err)
			}
			shutdown = true

		case r := <-q.sendMessageChan:
			queue.Topic[r.TopicName] = append(queue.Topic[r.TopicName], r.Message)

		case r := <-q.retrieveMessageChan:
			if _, ok := queue.Topic[r.TopicName]; ok {
				q.retrieveMessageResponseCh <- RetrieveMessageResponse{
					Message: queue.Topic[r.TopicName][len(queue.Topic[r.TopicName])-1],
				}
			}
			q.retrieveMessageResponseCh <- RetrieveMessageResponse{}

		case r := <-q.deleteChan:
			queue.Topic[r.TopicName][len(queue.Topic[r.TopicName])-1] = Message{}
		}
	}
}

func (q *QueueServer) loadQueues(ctx context.Context) (*Queues, error) {
	queue, err := q.fetchQueues(ctx)
	if err != nil {
		return nil, err
	}

	topic := map[string][]Message{}

	for t, mm := range queue.Topic {
		for _, m := range mm {
			mg := Message{
				MessageID: m.MessageID,
				Data:      m.Data,
				CretedAt:  m.CretedAt,
				ExpiresAt: m.ExpiresAt,
			}
			topic[t] = append(topic[t], mg)
		}
	}

	return &Queues{
		Topic: topic,
	}, nil
}

func (q *QueueServer) saveQueues(ctx context.Context, queue *Queues) error {
	panic("not implemented")
}

// Shutdown gracefully shutdown the Queue Service
func (q *QueueServer) Shutdown(ctx context.Context) error {
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
