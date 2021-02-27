package storage

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //to be used indirectly by mysql driver
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const mysqlDriver = "mysql"

// DatabaseIF is an interface for messaging queue
type DatabaseIF interface {
	Connect() error
	Test() error
	FetchAllTopics(ctx context.Context, id int) (*[]string, error)
	InsertPublisher(ctx context.Context, publisherID int, topicName string) error
	UpdateTopicIDIntoPublisher(ctx context.Context, publisherID int, topicID string) error
	RemoveTopicIDFromPublisher(ctx context.Context, publisherID int) error
	FetchQueues(ctx context.Context) (*Queue, error)
	GetTopicIDFromPublisher(ctx context.Context, publisherID int) (string, bool, error)
	GetTopicIDFromTopic(ctx context.Context, topicName string) (string, error)
	InsertMessageIntoMessage(ctx context.Context, publisherID int, topicID string, message Message) error
	GetSubscribedTopics(ctx context.Context, subscriberID int) ([]string, error)
	InsertSubscriberIDIntoSubscriber(ctx context.Context, subscriberID int) error
	InsertIntoSubscriberTopicMap(ctx context.Context, subscriberID int, topicID string) error
	RemoveTopicIDFromSubscriberTopicMap(ctx context.Context, subscriberID int, topicID string) error
	SaveQueues(ctx context.Context, livequeue *[]StoreQueue, deadQueue *[]StoreQueue) error
	RemoveMessagesFromQueue(ctx context.Context) error
}

// MysqlDB is the reciever type for DatabaseIF
type MysqlDB struct {
	Dsn string
	Cxn *sql.DB
	Log *logrus.Logger
}

// NewMysqlDB creates a new DatabaseIF for mysql db
func NewMysqlDB(dataSourseName string, log *logrus.Logger) (DatabaseIF, error) {
	mysql := &MysqlDB{
		Dsn: dataSourseName,
		Log: log,
	}

	if err := mysql.Connect(); err != nil {
		return nil, err
	}

	return mysql, nil
}

// Connect connects to the DB instance
func (m *MysqlDB) Connect() error {
	cxn, err := sql.Open(mysqlDriver, m.Dsn)
	if err != nil {
		return errors.Wrap(err, "could not create connection to MYSQL db")
	}

	m.Cxn = cxn
	if err := m.Test(); err != nil {
		return err
	}

	return nil
}

// Test test the DB connection to see if connection if necessary
func (m *MysqlDB) Test() error {
	if err := m.Cxn.Ping(); err != nil {
		return errors.Wrap(err, "could not conect to MYSQL db")
	}
	return nil
}

// GetTopicIDFromPublisher gets topicId from Publisher table based on given publisherId
func (m *MysqlDB) GetTopicIDFromPublisher(ctx context.Context, publisherID int) (string, bool, error) {
	var (
		topicID  string
		notFound bool
	)

	stmt := `SELECT IFNULL(topicId,"") FROM Publisher WHERE publisherId = ?`

	err := m.Cxn.QueryRowContext(ctx, stmt, publisherID).Scan(&topicID)
	if err != nil && err != sql.ErrNoRows {
		return "", notFound, err
	}

	if err == sql.ErrNoRows {
		notFound = true
		return "", notFound, nil
	}

	return topicID, notFound, nil
}

// UpdateTopicIDIntoPublisher updates topicId into Publisher Table
func (m *MysqlDB) UpdateTopicIDIntoPublisher(ctx context.Context, publisherID int, topicID string) error {
	stmt := `UPDATE Publisher SET topicId = ? WHERE publisherId = ?`

	_, err := m.Cxn.ExecContext(ctx, stmt, topicID, publisherID)
	if err != nil {
		return err
	}

	return nil
}

// GetTopicIDFromTopic gets topicId from Topic table based on given topicName
func (m *MysqlDB) GetTopicIDFromTopic(ctx context.Context, topicName string) (string, error) {
	var topicID string

	stmt := `SELECT topicId FROM MessagingQueue.Topic where name = ?`

	err := m.Cxn.QueryRowContext(ctx, stmt, topicName).Scan(&topicID)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	return topicID, nil
}

// RemoveTopicIDFromPublisher  removes topicId from Publisher based on given publisherId
func (m *MysqlDB) RemoveTopicIDFromPublisher(ctx context.Context, publisherID int) error {
	stmt := `UPDATE Publisher SET topicId = NULL WHERE publisherId = ?`

	_, err := m.Cxn.ExecContext(ctx, stmt, publisherID)
	if err != nil {
		return err
	}

	return nil
}

// FetchQueues fetches messages for the queue
func (m *MysqlDB) FetchQueues(ctx context.Context) (*Queue, error) {
	stmt := `SELECT Q.topicId,M.messageId,M.data,M.createdAt,M.expiredAt 
				FROM Queue as Q JOIN Message as M 
				ON Q.messageId = M.messageId`

	row, err := m.Cxn.QueryContext(ctx, stmt)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	defer row.Close()

	t := map[string][]Message{}

	for row.Next() {
		var topicID string
		m := Message{}

		if err := row.Scan(&topicID, &m.MessageID, &m.Data, &m.CretedAt, &m.ExpiresAt); err != nil {
			return nil, err
		}
		t[topicID] = append(t[topicID], m)
	}

	return &Queue{Topic: t}, nil
}

// FetchAllTopics fetches all the topics from Topic table
func (m *MysqlDB) FetchAllTopics(ctx context.Context, id int) (*[]string, error) {
	var topics []string

	stmt := `SELECT name FROM Topic`

	row, err := m.Cxn.QueryContext(ctx, stmt)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	for row.Next() {
		var topic string
		if err := row.Scan(&topic); err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}

	return &topics, nil
}

// InsertPublisher insert new Publisher to Publisher table
func (m *MysqlDB) InsertPublisher(ctx context.Context, publisherID int, topicName string) error {
	stmt := `INSERT INTO Publisher (publisherId, topicId) VALUES (?,?)`

	_, err := m.Cxn.Exec(stmt, publisherID, topicName)
	if err != nil {
		return err
	}

	return nil
}

// InsertMessageIntoMessage persists message info into Message table
func (m *MysqlDB) InsertMessageIntoMessage(ctx context.Context, publisherID int, topicID string, message Message) error {
	stmt := `INSERT INTO Message (messageId,data,createdAt,expiredAt,pubId,topicId) VALUES (?,?,?,?,?,?)`

	_, err := m.Cxn.ExecContext(ctx, stmt, message.MessageID, message.Data, message.CretedAt, message.ExpiresAt, publisherID, topicID)
	if err != nil {
		return err
	}

	return nil
}

// GetSubscribedTopics fetches all the topics subscrebed by client
func (m *MysqlDB) GetSubscribedTopics(ctx context.Context, subscriberID int) ([]string, error) {
	var topics []string

	stmt := `SELECT T.name FROM Topic as T 
				RIGHT JOIN SubscriberTopicMap AS S
				ON T.topicId = S.topicId WHERE S.subscriberId = ?`

	row, err := m.Cxn.QueryContext(ctx, stmt, subscriberID)
	if err != nil {
		return []string{}, err
	}

	for row.Next() {
		var topic string
		if err := row.Scan(&topic); err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}

	return topics, nil
}

// InsertSubscriberIDIntoSubscriber inserts new subscriber into Subscriber table
func (m *MysqlDB) InsertSubscriberIDIntoSubscriber(ctx context.Context, subscriberID int) error {
	var id int

	stmt := `SELECT subscriberId FROM Subscriber WHERE subscriberId = ?`

	err := m.Cxn.QueryRowContext(ctx, stmt, subscriberID).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if id == 0 {
		stmt = `INSERT INTO Subscriber (subscriberId) VALUES (?)`

		_, err = m.Cxn.ExecContext(ctx, stmt, subscriberID)
		if err != nil {
			return err
		}
	}

	return nil
}

// InsertIntoSubscriberTopicMap inserts susbscriber to topic mapping
func (m *MysqlDB) InsertIntoSubscriberTopicMap(ctx context.Context, subscriberID int, topicID string) error {
	stmt := `INSERT INTO SubscriberTopicMap (subscriberId,topicId) VALUES (?,?)`

	_, err := m.Cxn.ExecContext(ctx, stmt, subscriberID, topicID)
	if err != nil {
		return err
	}

	return nil
}

// RemoveTopicIDFromSubscriberTopicMap remove topicId andm subscriberId from SubscriberTopicMap table
func (m *MysqlDB) RemoveTopicIDFromSubscriberTopicMap(ctx context.Context, subscriberID int, topicID string) error {
	stmt := `DELETE FROM SubscriberTopicMap WHERE subscriberId = ? AND topicId = ?`

	_, err := m.Cxn.ExecContext(ctx, stmt, subscriberID, topicID)
	if err != nil {
		return err
	}

	return nil
}

// SaveQueues persists all the message from topic to db
func (m *MysqlDB) SaveQueues(ctx context.Context, livequeue *[]StoreQueue, deadQueue *[]StoreQueue) error {
	stmt := `INSERT INTO Queue (queueId,topicId,messageId) VALUES (?,?,?)`

	for _, qs := range *livequeue {
		_, err := m.Cxn.ExecContext(ctx, stmt, qs.QueuID, qs.TopicID, qs.MessageID)
		if err != nil {
			return err
		}
	}

	dlqStmt := `INSERT INTO DLQ (dlqId,topicId,messageId) VALUES (?,?,?)`

	for _, qs := range *deadQueue {
		_, err := m.Cxn.ExecContext(ctx, dlqStmt, qs.QueuID, qs.TopicID, qs.MessageID)
		if err != nil {
			return err
		}
	}

	return nil
}

// RemoveMessagesFromQueue clear Queue table
func (m *MysqlDB) RemoveMessagesFromQueue(ctx context.Context) error {
	stmt := `DELETE FROM Queue`

	_, err := m.Cxn.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}

	return nil
}
