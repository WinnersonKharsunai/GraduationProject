package queueserver

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //...
	"github.com/pkg/errors"
)

const mysqlDriver = "mysql"

func (q *QueueServer) connect() error {

	cxn, err := sql.Open(mysqlDriver, q.Dsn)
	if err != nil {
		return errors.Wrap(err, "could not create connection to MYSQL db")
	}

	q.Cxn = cxn

	if err := q.Cxn.Ping(); err != nil {
		return errors.Wrap(err, "could not conect to MYSQL db")
	}

	return nil
}

func (q *QueueServer) fetchQueues(ctx context.Context) (*Queue, error) {

	stmt := `SELECT M.topicName,M.Data,M.createdAt,M.expiredAt 
			 FROM Topic AS T RIGHT JOIN Message AS M 
			 ON T.name = M.topicName`

	row, err := q.Cxn.QueryContext(ctx, stmt)
	if err != nil && err == sql.ErrNoRows {
		return nil, err
	}

	defer row.Close()

	t := map[string][]Message{}

	for row.Next() {
		var topicName string
		m := Message{}

		if err := row.Scan(&topicName, &m.Data, &m.CretedAt, &m.ExpiresAt); err != nil {
			return nil, err
		}

		t[topicName] = append(t[topicName], m)
	}

	return &Queue{Topic: t}, nil
}
