package task

import (
	"context"
	"database/sql"
	"fmt"
)

type InsertDB struct {
	DB   *sql.DB
	Data string
}

func (t *InsertDB) Process(ctx context.Context) (string, error) {
	_, err := t.DB.ExecContext(ctx, "INSERT INTO users (name) VALUES (?)", t.Data)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Insertef: %s", t.Data), nil
}
