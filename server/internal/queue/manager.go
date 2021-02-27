package queue

import (
	"context"
	"fmt"
)

func (q *Queue) queueService(queue Queues) {

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
			queue.Topic[r.TopicID] = append(queue.Topic[r.TopicID], r.Message)

		case r := <-q.retrieveMessageChan:
			if _, ok := queue.Topic[r.TopicName]; ok {
				q.retrieveMessageResponseCh <- RetrieveMessageResponse{
					Message: queue.Topic[r.TopicName][len(queue.Topic[r.TopicName])-1],
				}
			}

		case r := <-q.deleteChan:
			queue.Topic[r.TopicName][len(queue.Topic[r.TopicName])-1] = Message{}
		}
	}
}

func (q *Queue) loadQueues(ctx context.Context) (*Queues, error) {
	queue, err := q.db.FetchQueues(ctx)
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

func (q *Queue) saveQueues(ctx context.Context, queue *Queues) error {
	panic("not implemented")
}
