package storage_test

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/WinnersonKharsunai/GraduationProject/server/internal/storage"
)

func TestGetTopicIDFromPublisher_Fail(t *testing.T) {
	publisherID := 5000
	expectedErr := errors.New("failed to get publisherId")

	mock, db := mysqlMock()
	stmt := `SELECT IFNULL\(topicId,""\) FROM Publisher WHERE publisherId = \?`
	mock.ExpectQuery(stmt).WillReturnError(expectedErr)

	_, _, err := db.GetTopicIDFromPublisher(context.Background(), publisherID)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
}

func TestGetTopicIDFromPublisher_NotFoundFail(t *testing.T) {
	publisherID := 5000

	mock, db := mysqlMock()
	stmt := `SELECT IFNULL\(topicId,""\) FROM Publisher WHERE publisherId = \?`
	mock.ExpectQuery(stmt).WillReturnError(sql.ErrNoRows)

	_, notFound, err := db.GetTopicIDFromPublisher(context.Background(), publisherID)
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}

	if notFound != true {
		t.Fatalf("expected: notFound, got: %v", notFound)
	}
}

func TestGetTopicIDFromPublisher_Pass(t *testing.T) {
	publisherID := 5000
	mock, db := mysqlMock()

	columns := []string{"topicId"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow("12345")

	stmt := `SELECT IFNULL\(topicId,""\) FROM Publisher WHERE publisherId = \?`
	mock.ExpectQuery(stmt).WillReturnRows(rows)

	topicID, _, err := db.GetTopicIDFromPublisher(context.Background(), publisherID)
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}

	if topicID != "12345" {
		t.Fatalf("expected: topicId = '12345', got: %v", topicID)
	}
}

func TestUpdateTopicIDIntoPublisher_Fail(t *testing.T) {
	publisherID := 5000
	topicID := "12345"
	expectedErr := errors.New("failed to update")

	mock, db := mysqlMock()
	stmt := `UPDATE Publisher SET topicId = \? WHERE publisherId = \?`
	mock.ExpectExec(stmt).WillReturnError(expectedErr)

	err := db.UpdateTopicIDIntoPublisher(context.Background(), publisherID, topicID)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
}

func TestUpdateTopicIDIntoPublisher_Pass(t *testing.T) {
	publisherID := 5000
	topicID := "12345"

	mock, db := mysqlMock()
	stmt := `UPDATE Publisher SET topicId = \? WHERE publisherId = \?`
	mock.ExpectExec(stmt).WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.UpdateTopicIDIntoPublisher(context.Background(), publisherID, topicID)
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}
}

func TestGetTopicIDFromTopic_Fail(t *testing.T) {
	topicName := "Golang"
	expectedErr := errors.New("failed to get topicId")

	mock, db := mysqlMock()
	stmt := `SELECT topicId FROM MessagingQueue.Topic where name = \?`
	mock.ExpectQuery(stmt).WillReturnError(expectedErr)

	_, err := db.GetTopicIDFromTopic(context.Background(), topicName)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
}

func TestGetTopicIDFromTopic_Pass(t *testing.T) {
	topicName := "Golang"
	mock, db := mysqlMock()

	columns := []string{"topicId"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow("12345")

	stmt := `SELECT topicId FROM MessagingQueue.Topic where name = \?`
	mock.ExpectQuery(stmt).WillReturnRows(rows)

	topicID, err := db.GetTopicIDFromTopic(context.Background(), topicName)
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}

	if topicID != "12345" {
		t.Fatalf("expected: topicId = '12345', got: %v", topicID)
	}
}

func TestRemoveTopicIDFromPublisher_Fail(t *testing.T) {
	publisherID := 5000
	expectedErr := errors.New("failed to update")

	mock, db := mysqlMock()
	stmt := `UPDATE Publisher SET topicId = NULL WHERE publisherId = \?`
	mock.ExpectExec(stmt).WillReturnError(expectedErr)

	err := db.RemoveTopicIDFromPublisher(context.Background(), publisherID)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
}

func TestRemoveTopicIDFromPublisher_Pass(t *testing.T) {
	publisherID := 5000

	mock, db := mysqlMock()
	stmt := `UPDATE Publisher SET topicId = NULL WHERE publisherId = \?`
	mock.ExpectExec(stmt).WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.RemoveTopicIDFromPublisher(context.Background(), publisherID)
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}
}

func TestFetchQueues_Fail(t *testing.T) {
	expectedErr := errors.New("failed to fetch")

	mock, db := mysqlMock()
	stmt := `SELECT Q.topicId,M.messageId,M.data,M.createdAt,M.expiredAt 
				FROM Queue as Q JOIN Message as M 
				ON Q.messageId = M.messageId`
	mock.ExpectQuery(stmt).WillReturnError(expectedErr)

	_, err := db.FetchQueues(context.Background())
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
}

func TestFetchQueues_Pass(t *testing.T) {
	mock, db := mysqlMock()

	want := &storage.Queue{
		Topic: map[string][]storage.Message{
			"12345": {
				{
					MessageID: "123",
					Data:      "test",
					CretedAt:  "2021-02-27 20:03:09",
					ExpiresAt: "2021-02-27 20:04:09",
				},
			},
		},
	}

	columns := []string{"topicId", "messageId", "data", "createdAt", "expiredAt"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow("12345", "123", "test", "2021-02-27 20:03:09", "2021-02-27 20:04:09")

	stmt := `SELECT Q.topicId,M.messageId,M.data,M.createdAt,M.expiredAt 
				FROM Queue as Q JOIN Message as M 
				ON Q.messageId = M.messageId`
	mock.ExpectQuery(stmt).WillReturnRows(rows)

	queue, err := db.FetchQueues(context.Background())
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}

	if !reflect.DeepEqual(queue, want) {
		t.Fatalf("expected: %v, got: %v", want, queue)
	}
}

func TestFetchAllTopics_Fail(t *testing.T) {
	publisherID := 5000
	expectedErr := errors.New("failed to fetch")

	mock, db := mysqlMock()
	stmt := `SELECT name FROM Topic`
	mock.ExpectQuery(stmt).WillReturnError(expectedErr)

	_, err := db.FetchAllTopics(context.Background(), publisherID)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
}

func TestFetchAllTopics_Pass(t *testing.T) {
	publisherID := 5000
	mock, db := mysqlMock()

	columns := []string{"name"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow("golang")

	stmt := `SELECT name FROM Topic`
	mock.ExpectQuery(stmt).WillReturnRows(rows)

	_, err := db.FetchAllTopics(context.Background(), publisherID)
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}
}

func TestInsertPublisher_Fail(t *testing.T) {
	publisherID := 5000
	topicName := "Test"
	expectedErr := errors.New("failed to insert")

	mock, db := mysqlMock()
	stmt := `INSERT INTO Publisher \(publisherId, topicId\) VALUES \(\?,\?\)`
	mock.ExpectExec(stmt).WillReturnError(expectedErr)

	err := db.InsertPublisher(context.Background(), publisherID, topicName)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
}

func TestInsertPublisher_Pass(t *testing.T) {
	publisherID := 5000
	topicName := "Test"

	mock, db := mysqlMock()
	stmt := `INSERT INTO Publisher \(publisherId, topicId\) VALUES \(\?,\?\)`
	mock.ExpectExec(stmt).WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.InsertPublisher(context.Background(), publisherID, topicName)
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}
}

func TestInsertMessageIntoMessage_Fail(t *testing.T) {
	publisherID := 5000
	topicID := "12345"
	message := storage.Message{
		MessageID: "123",
		Data:      "test data",
		CretedAt:  "2021-02-27 20:03:09",
		ExpiresAt: "2021-02-27 20:04:09",
	}
	expectedErr := errors.New("failed to insert")

	mock, db := mysqlMock()
	stmt := `INSERT INTO Message \(messageId,data,createdAt,expiredAt,pubId,topicId\) VALUES \(\?,\?,\?,\?,\?,\?\)`
	mock.ExpectExec(stmt).WillReturnError(expectedErr)

	err := db.InsertMessageIntoMessage(context.Background(), publisherID, topicID, message)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
}

func TestInsertMessageIntoMessage_Pass(t *testing.T) {
	publisherID := 5000
	topicID := "12345"
	message := storage.Message{
		MessageID: "123",
		Data:      "test data",
		CretedAt:  "2021-02-27 20:03:09",
		ExpiresAt: "2021-02-27 20:04:09",
	}

	mock, db := mysqlMock()
	stmt := `INSERT INTO Message \(messageId,data,createdAt,expiredAt,pubId,topicId\) VALUES \(\?,\?,\?,\?,\?,\?\)`
	mock.ExpectExec(stmt).WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.InsertMessageIntoMessage(context.Background(), publisherID, topicID, message)
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}
}

func TestGetSubscribedTopics_Fail(t *testing.T) {
	subscriberID := 6000
	expectedErr := errors.New("failed to get topics")

	mock, db := mysqlMock()
	stmt := `SELECT T.name FROM Topic as T 
				RIGHT JOIN SubscriberTopicMap AS S
				ON T.topicId = S.topicId WHERE S.subscriberId = \?`
	mock.ExpectQuery(stmt).WillReturnError(expectedErr)

	_, err := db.GetSubscribedTopics(context.Background(), subscriberID)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
}

func TestGetSubscribedTopics_Pass(t *testing.T) {
	subscriberID := 6000
	mock, db := mysqlMock()

	columns := []string{"name"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow("golang")

	stmt := `SELECT T.name FROM Topic as T 
				RIGHT JOIN SubscriberTopicMap AS S
				ON T.topicId = S.topicId WHERE S.subscriberId = \?`
	mock.ExpectQuery(stmt).WillReturnRows(rows)

	_, err := db.GetSubscribedTopics(context.Background(), subscriberID)
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}
}

func TestInsertSubscriberIDIntoSubscriber_FailedToGetSubscriberID(t *testing.T) {
	subscriberID := 6000
	expectedErr := errors.New("failed to  get subscriberId")

	mock, db := mysqlMock()
	stmt := `SELECT subscriberId FROM Subscriber WHERE subscriberId = \?`
	mock.ExpectQuery(stmt).WillReturnError(expectedErr)

	err := db.InsertSubscriberIDIntoSubscriber(context.Background(), subscriberID)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
}

func TestInsertSubscriberIDIntoSubscriber_FailedToInsert(t *testing.T) {
	subscriberID := 6000
	expectedErr := errors.New("failed to insert subscriberId")
	mock, db := mysqlMock()

	columns := []string{"subscriberId"}
	rows := sqlmock.NewRows(columns)

	stmt := `SELECT subscriberId FROM Subscriber WHERE subscriberId = \?`
	mock.ExpectQuery(stmt).WillReturnRows(rows)

	stmt = `INSERT INTO Subscriber \(subscriberId\) VALUES \(\?\)`
	mock.ExpectExec(stmt).WillReturnError(expectedErr)

	err := db.InsertSubscriberIDIntoSubscriber(context.Background(), subscriberID)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
}

func TestInsertSubscriberIDIntoSubscriber_Pass(t *testing.T) {
	subscriberID := 6000
	mock, db := mysqlMock()

	columns := []string{"subscriberId"}
	rows := sqlmock.NewRows(columns)

	stmt := `SELECT subscriberId FROM Subscriber WHERE subscriberId = \?`
	mock.ExpectQuery(stmt).WillReturnRows(rows)

	stmt = `INSERT INTO Subscriber \(subscriberId\) VALUES \(\?\)`
	mock.ExpectExec(stmt).WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.InsertSubscriberIDIntoSubscriber(context.Background(), subscriberID)
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}
}

func TestInsertIntoSubscriberTopicMap_Fail(t *testing.T) {
	subscriberID := 6000
	topicID := "12345"
	expectedErr := errors.New("failed to insert")

	mock, db := mysqlMock()
	stmt := `INSERT INTO SubscriberTopicMap \(subscriberId,topicId\) VALUES \(\?,\?\)`
	mock.ExpectExec(stmt).WillReturnError(expectedErr)

	err := db.InsertIntoSubscriberTopicMap(context.Background(), subscriberID, topicID)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
}

func TestInsertIntoSubscriberTopicMap_Pass(t *testing.T) {
	subscriberID := 6000
	topicID := "12345"

	mock, db := mysqlMock()
	stmt := `INSERT INTO SubscriberTopicMap \(subscriberId,topicId\) VALUES \(\?,\?\)`
	mock.ExpectExec(stmt).WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.InsertIntoSubscriberTopicMap(context.Background(), subscriberID, topicID)
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}
}

func TestRemoveTopicIDFromSubscriberTopicMap_Fail(t *testing.T) {
	subscriberID := 6000
	topicID := "12345"
	expectedErr := errors.New("failed to delete")

	mock, db := mysqlMock()
	stmt := `DELETE FROM SubscriberTopicMap WHERE subscriberId = \? AND topicId = \?`
	mock.ExpectExec(stmt).WillReturnError(expectedErr)

	err := db.RemoveTopicIDFromSubscriberTopicMap(context.Background(), subscriberID, topicID)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
}

func TestRemoveTopicIDFromSubscriberTopicMap_Pass(t *testing.T) {
	subscriberID := 6000
	topicID := "12345"

	mock, db := mysqlMock()
	stmt := `DELETE FROM SubscriberTopicMap WHERE subscriberId = \? AND topicId = \?`
	mock.ExpectExec(stmt).WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.RemoveTopicIDFromSubscriberTopicMap(context.Background(), subscriberID, topicID)
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}
}

func TestSaveQueues_FailedToInsertToQueue(t *testing.T) {
	liveQueue := &[]storage.StoreQueue{
		{
			QueuID:    "queue",
			TopicID:   "12334",
			MessageID: "message123",
		},
	}

	expectedErr := errors.New("failed to insert to Queue")

	mock, db := mysqlMock()
	stmt := `INSERT INTO Queue \(queueId,topicId,messageId\) VALUES \(\?,\?,\?\)`
	mock.ExpectExec(stmt).WillReturnError(expectedErr)

	err := db.SaveQueues(context.Background(), liveQueue, true)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
}

func TestSaveQueues_FailedToInsertToDLQ(t *testing.T) {
	deadQueue := &[]storage.StoreQueue{
		{
			QueuID:    "queue",
			TopicID:   "12334",
			MessageID: "message123",
		},
	}

	expectedErr := errors.New("failed to insert to Queue")

	mock, db := mysqlMock()
	dlqStmt := `INSERT INTO DLQ \(dlqId,topicId,messageId\) VALUES \(\?,\?,\?\)`
	mock.ExpectExec(dlqStmt).WillReturnError(expectedErr)

	err := db.SaveQueues(context.Background(), deadQueue, false)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
}

func TestSaveQueues_Pass(t *testing.T) {
	liveQueue := &[]storage.StoreQueue{
		{
			QueuID:    "queue",
			TopicID:   "12334",
			MessageID: "message123",
		},
	}

	mock, db := mysqlMock()
	stmt := `INSERT INTO Queue \(queueId,topicId,messageId\) VALUES \(\?,\?,\?\)`
	mock.ExpectExec(stmt).WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.SaveQueues(context.Background(), liveQueue, true)
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}
}

func TestRemoveMessagesFromQueue_Fail(t *testing.T) {
	expectedErr := errors.New("failed to delete")

	mock, db := mysqlMock()
	stmt := `DELETE FROM Queue`
	mock.ExpectExec(stmt).WillReturnError(expectedErr)

	err := db.RemoveMessagesFromQueue(context.Background())
	if err.Error() != expectedErr.Error() {
		t.Fatalf("expected: %v, got: %v", expectedErr, err)
	}
}

func TestRemoveMessagesFromQueue_Pass(t *testing.T) {

	mock, db := mysqlMock()
	stmt := `DELETE FROM Queue`
	mock.ExpectExec(stmt).WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.RemoveMessagesFromQueue(context.Background())
	if err != nil {
		t.Fatalf("expected: nil, got: %v", err)
	}
}

func mysqlMock() (sqlmock.Sqlmock, storage.MysqlDB) {
	dbCxn, mock, _ := sqlmock.New()
	db := storage.MysqlDB{
		Cxn: dbCxn,
	}
	return mock, db
}
