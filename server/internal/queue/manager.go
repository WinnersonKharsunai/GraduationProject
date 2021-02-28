package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/WinnersonKharsunai/GraduationProject/server/internal/storage"
	"github.com/google/uuid"
)

func (q *Queue) queueService(queue Queues) {
	shutdown := false
	for !shutdown {
		liveMesage := queue.Topic
		deadMessage := queue.DLQ
		fmt.Println("\nQueue Topic:", liveMesage)
		fmt.Println("\nDLQ:", deadMessage)

		select {
		case <-q.shutdownChan:
			if err := q.saveQueues(context.Background(), &Queues{Topic: liveMesage, DLQ: deadMessage}); err != nil {
				fmt.Println(err)
			}
			shutdown = true

		case r := <-q.sendMessageChan:
			liveMesage[r.TopicID] = append(liveMesage[r.TopicID], r.Message)

		case r := <-q.retrieveMessageChan:
			if msg, ok := liveMesage[r.TopicID]; ok {

				length := len(msg)
				if length > 0 {
					var i int
					m := Message{}
					for i, m = range msg {
						expTime := getEpochTime(m.ExpiresAt)
						curTime := getCurrentTime()

						if curTime >= expTime {
							deadMessage[r.TopicID] = append(deadMessage[r.TopicID], m)
							if i == length-1 {
								liveMesage[r.TopicID] = []Message{}
							} else {
								fmt.Println("i ws here last", i, length, liveMesage[r.TopicID][i+1:])
								liveMesage[r.TopicID] = liveMesage[r.TopicID][i+1:]
							}
						} else {
							q.retrieveMessageResponseCh <- RetrieveMessageResponse{Message: m}
							break
						}
					}

					if len(msg) == i+1 {
						q.retrieveMessageResponseCh <- RetrieveMessageResponse{}
					}

				} else {
					q.retrieveMessageResponseCh <- RetrieveMessageResponse{}
				}
			}
		}
	}
	q.processWg.Done()
}

func (q *Queue) loadQueues(ctx context.Context) (*Queues, error) {
	queue, err := q.db.FetchQueues(ctx)
	if err != nil {
		return nil, err
	}

	topic := map[string][]Message{}
	dlq := map[string][]Message{}

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
		DLQ:   dlq,
	}, nil
}

func (q *Queue) clearQueue(ctx context.Context) error {
	err := q.db.RemoveMessagesFromQueue(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (q *Queue) saveQueues(ctx context.Context, queue *Queues) error {
	liveMsgs := []storage.StoreQueue{}
	deadMsgs := []storage.StoreQueue{}

	for topicID, msg := range queue.Topic {
		for _, m := range msg {
			liveMsg := storage.StoreQueue{
				QueuID:    uuid.New().String(),
				MessageID: m.MessageID,
				TopicID:   topicID,
			}
			liveMsgs = append(liveMsgs, liveMsg)
		}
	}

	for topicID, msg := range queue.DLQ {
		for _, m := range msg {
			deadMsg := storage.StoreQueue{
				QueuID:    uuid.New().String(),
				MessageID: m.MessageID,
				TopicID:   topicID,
			}
			deadMsgs = append(liveMsgs, deadMsg)
		}
	}

	if err := q.db.SaveQueues(ctx, &liveMsgs, &deadMsgs); err != nil {
		return err
	}

	return nil
}

func getEpochTime(t string) int64 {
	thetime, _ := time.Parse("2006-01-02 15:04:05", t)
	return thetime.Unix()
}

func getCurrentTime() int64 {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	thetime, _ := time.Parse("2006-01-02 15:04:05", now)
	return thetime.Unix()
}
