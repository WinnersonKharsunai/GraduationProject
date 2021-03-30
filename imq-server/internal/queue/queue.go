package queue

import (
	"context"
	"errors"
<<<<<<< HEAD
	"time"
	"fmt"
=======
	"fmt"
	"time"
>>>>>>> 9fe39465475b121a78fe3f5e4b7a5638b6c0a469

	"github.com/WinnersonKharsunai/GraduationProject/server/internal/storage"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// ImqQueueIF is the inteerface for the Queue
type ImqQueueIF interface {
	SendMessage(ctx context.Context, message SendMessageRequest) error
	RetrieveMessage(ctx context.Context, topicID string) (*Message, error)
	BackUpQueue(ctx context.Context) error
<<<<<<< HEAD
	loadQueue() error
	saveQueue(ctx context.Context) error
	clearQueueFromDb() error
=======
>>>>>>> 9fe39465475b121a78fe3f5e4b7a5638b6c0a469
}

// Queue is the concrete implementztion for Queue
type Queue struct {
	log       *logrus.Logger
	db        storage.DatabaseIF
	LiveQueue map[string][]Message
	DeadQueue map[string][]Message
}

// NewQueue is the factory function for the Queue
func NewQueue(log *logrus.Logger, db storage.DatabaseIF) (ImqQueueIF, error) {
	q := &Queue{
		log:       log,
		db:        db,
		LiveQueue: map[string][]Message{},
		DeadQueue: map[string][]Message{},
	}

	if err := q.loadQueue(); err != nil {
<<<<<<< HEAD
		q.log.Errorf("failed to load queue: %v", err)
		return nil, err
	}

	q.log.Infof("QUEUE: %v", q.LiveQueue)

	if err := q.clearQueueFromDb(); err != nil {
		q.log.Errorf("failed to clear queue table: %v", err)
=======
		return nil, err
	}

	if err := q.clearQueueFromDb(); err != nil {
>>>>>>> 9fe39465475b121a78fe3f5e4b7a5638b6c0a469
		return nil, err
	}

	return q, nil
}

// SendMessage push message to the queue
func (q *Queue) SendMessage(ctx context.Context, request SendMessageRequest) error {
	if request.TopicID == "" {
<<<<<<< HEAD
		return errors.New("you are not register to any topics")
=======
		return errors.New("topicId cannot be empty")
>>>>>>> 9fe39465475b121a78fe3f5e4b7a5638b6c0a469
	}

	if request.Message.MessageID == "" || request.Message.Data == "" {
		return errors.New("message cannot be empty")
	}
<<<<<<< HEAD
	
	length := len(request.Message.Data) 

	fmt.Println(length)

	q.LiveQueue[request.TopicID] = append(q.LiveQueue[request.TopicID], request.Message)

	q.log.Infof("QUEUE: %v", q.LiveQueue)
=======

	q.LiveQueue[request.TopicID] = append(q.LiveQueue[request.TopicID], request.Message)

	fmt.Println(q.LiveQueue)
>>>>>>> 9fe39465475b121a78fe3f5e4b7a5638b6c0a469

	return nil
}

// RetrieveMessage pull message from the queue
func (q *Queue) RetrieveMessage(ctx context.Context, topicID string) (*Message, error) {
	for {
<<<<<<< HEAD
		msg, err := peekMessage(q.LiveQueue[topicID])
		if err != nil {
			q.log.Infof("QUEUE: %v", q.LiveQueue)
			q.log.Infof("DLQ: %v", q.DeadQueue)
=======
		fmt.Println(q.LiveQueue[topicID])
		msg, err := peekMessage(q.LiveQueue[topicID])
		if err != nil {
>>>>>>> 9fe39465475b121a78fe3f5e4b7a5638b6c0a469
			return nil, err
		}

		if isExpired(msg.ExpiresAt) {
			pushToDeadMessage(msg, q.DeadQueue[topicID])
			copy(q.LiveQueue[topicID][0:], q.LiveQueue[topicID][1:])
			q.LiveQueue[topicID] = q.LiveQueue[topicID][:len(q.LiveQueue[topicID])-1]
		} else {
<<<<<<< HEAD
			q.log.Infof("QUEUE: %v", q.LiveQueue)
			q.log.Infof("DLQ: %v", q.DeadQueue)
=======
>>>>>>> 9fe39465475b121a78fe3f5e4b7a5638b6c0a469
			return &msg, nil
		}
	}
}

// BackUpQueue store the data from queue to db
func (q *Queue) BackUpQueue(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		if err := q.saveQueue(ctx); err != nil {
			q.log.Errorf("failed to backup queue: %v", err)
		}
<<<<<<< HEAD

		q.log.Infof("queue has been successfully backed up")

=======
>>>>>>> 9fe39465475b121a78fe3f5e4b7a5638b6c0a469
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *Queue) loadQueue() error {
	liveQueue, err := q.db.FetchQueues(context.Background())
	if err != nil {
<<<<<<< HEAD
=======
		fmt.Println(err)
>>>>>>> 9fe39465475b121a78fe3f5e4b7a5638b6c0a469
		return nil
	}

	if len(liveQueue.Topic) > 0 {
		for k, m := range liveQueue.Topic {
			for _, msg := range m {
				mm := Message{
					MessageID: msg.MessageID,
					Data:      msg.Data,
					CretedAt:  msg.CretedAt,
					ExpiresAt: msg.ExpiresAt,
				}
				q.LiveQueue[k] = append(q.LiveQueue[k], mm)
			}
		}
	}

	return nil
}

func (q *Queue) saveQueue(ctx context.Context) error {
	liveQueueData := getQueueData(q.LiveQueue)
	if err := q.db.SaveQueues(ctx, &liveQueueData, true); err != nil {
		return err
	}

	deadQueueData := getQueueData(q.LiveQueue)
	if err := q.db.SaveQueues(ctx, &deadQueueData, false); err != nil {
		return err
	}

	return nil
}

<<<<<<< HEAD
func (q *Queue) clearQueueFromDb() error {
	err := q.db.RemoveMessagesFromQueue(context.Background())
	if err != nil {
		return err
	}

	return nil
}

=======
>>>>>>> 9fe39465475b121a78fe3f5e4b7a5638b6c0a469
func getQueueData(queue map[string][]Message) []storage.StoreQueue {
	data := []storage.StoreQueue{}
	for topicID, messages := range queue {
		for _, m := range messages {
			msg := storage.StoreQueue{
				QueuID:    uuid.New().String(),
				TopicID:   topicID,
				MessageID: m.MessageID,
			}
			data = append(data, msg)
		}
	}
	return data
}

func pushToDeadMessage(msg Message, msgs []Message) {
	msgs = append(msgs, msg)
}

func peekMessage(msg []Message) (Message, error) {
	if len(msg) <= 0 {
		return Message{}, errors.New("no message present in queue")
	}
	return msg[0], nil
}

func isExpired(t string) bool {
	expTime := getEpochTime(t)
	curTime := getCurrentTime()
	if curTime > expTime {
		return true
	}
	return false
}

<<<<<<< HEAD
=======
func (q *Queue) clearQueueFromDb() error {
	err := q.db.RemoveMessagesFromQueue(context.Background())
	if err != nil {
		return err
	}

	return nil
}

>>>>>>> 9fe39465475b121a78fe3f5e4b7a5638b6c0a469
func getEpochTime(t string) int64 {
	thetime, _ := time.Parse("2006-01-02 15:04:05", t)
	return thetime.Unix()
}

func getCurrentTime() int64 {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	thetime, _ := time.Parse("2006-01-02 15:04:05", now)
	return thetime.Unix()
}
