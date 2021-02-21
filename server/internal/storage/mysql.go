package storage

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //...
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var mysqlDriver = "mysql"

// DatabaseIF is an interface for messaging queue
type DatabaseIF interface {
	Connect() error
	Test() error
	FetchAllTopics(ctx context.Context) ([]string, error)
	InsertPublisher(ctx context.Context, id int, topicName string) error
	RemovePublisher(ctx context.Context, id int) error
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

// FetchAllTopics ...
func (m *MysqlDB) FetchAllTopics(ctx context.Context) ([]string, error) {
	return []string{}, nil
}

// InsertPublisher ...
func (m *MysqlDB) InsertPublisher(ctx context.Context, id int, topicName string) error {
	return nil
}

// RemovePublisher  ...
func (m *MysqlDB) RemovePublisher(ctx context.Context, id int) error {
	return nil
}
