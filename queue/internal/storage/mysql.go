package storage

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var mysqlDriver = "mysql"

// DatabaseIF is an interface for messaging queue
type DatabaseIF interface {
	Connect() error
	Test() error
	CreateTopic(ctx context.Context, tp *NewTopicInfo) error
	DeleteTopic(ctx context.Context, topicName string) error
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
func (mysql *MysqlDB) Connect() error {

	cxn, err := sql.Open(mysqlDriver, mysql.Dsn)
	if err != nil {
		return errors.Wrap(err, "could not create connection to MYSQL db")
	}
	mysql.Cxn = cxn
	if err := mysql.Test(); err != nil {
		return err
	}

	return nil
}

// Test test the DB connection to see if connection if necessary
func (mysql *MysqlDB) Test() error {

	if err := mysql.Cxn.Ping(); err != nil {
		return errors.Wrap(err, "could not conect to MYSQL db")
	}
	return nil
}

// CreateTopic ...
func (mysql *MysqlDB) CreateTopic(ctx context.Context, tp *NewTopicInfo) error {
	return nil
}

// DeleteTopic ..
func (mysql *MysqlDB) DeleteTopic(ctx context.Context, topicName string) error {
	return nil
}
