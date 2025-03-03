package postgres

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	dbDriver = "postgres"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(user, password, host, port, dbName string) (*Postgres, error) {
	var err error

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=enabled", host, port, user, password, dbName)
	db, err := sql.Open(dbDriver, connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Printf("postgres ping failure: %v\n", err)
		return nil, err
	}
	return &Postgres{db: db}, nil
}

func (p *Postgres) Close() {
	err := p.db.Close()
	if err != nil {
		log.Printf("postgres close failure: %v\n", err)
	}
}

func (p Postgres) InsertUser(userName string) error {
	p.db.Exec("INSERT...")
	return nil
}

func (p Postgres) SelectSingleUser(userName string) (string, error) {
	p.db.Exec("SELECT...")
	return "user", nil
}
