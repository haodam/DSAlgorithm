package task

import (
	"context"
	"database/sql"
)

type InsertDB struct {
	DB   *sql.DB
	Data []byte
}

func (t *InsertDB) Process(ctx context.Context) (string, error) {
	return "", nil
}
