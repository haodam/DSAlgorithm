package main

import (
	"github.com/haodam/DSAlgorithm/pkg/interface/postgres"
	"log"
	"os"
)

type dbContract interface {
	Close()
	InsertUser(username string) error
	SelectSingleUser(username string) (string, error)
}

type Application struct {
	db dbContract
}

func (app *Application) Run() {
	userName := "user1"
	err := app.db.InsertUser(userName)
	if err != nil {
		log.Printf("couldn't insert user: %s\n", userName)
	}

	user, err := app.db.SelectSingleUser(userName)
	if err != nil {
		log.Printf("couldn't find user: %s\n", userName)
	}
	log.Printf("user is %s\n", user)
}

func NewApplication(db dbContract) *Application {
	return &Application{
		db: db,
	}
}

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	db, err := postgres.NewPostgres(dbUser, dbPassword, dbHost, dbPort, dbName)
	if err != nil {
		log.Fatalf("failed to initiate dbase connection: %v", err)
	}

	defer db.Close()
	app := NewApplication(db)
	app.Run()
}
