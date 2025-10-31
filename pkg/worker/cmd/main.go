package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/haodam/DSAlgorithm/pkg/worker/task"
	"github.com/haodam/DSAlgorithm/pkg/worker/worker_pool"
)

func main() {

	ctx := context.Background()

	// Connect DB
	db, err := sql.Open("mysql", "xxxxxxxxxxxxxx")
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	// create a worker pool
	worker := worker_pool.NewWorkerPool[string](ctx,
		worker_pool.WithMaxWorkers(20),
		worker_pool.WithTimeout(10*time.Minute),
	)

	worker.Start()

	// Read file
	file, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(file), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tasks := task.InsertDB{
			DB:   db,
			Data: line,
		}
		if err := worker.Submit(&tasks); err != nil {
			log.Println("submit error:", err)
			break
		}
	}

	go func() {
		worker.Shutdown()
	}()

	results, errors := worker.CollectResults()

	fmt.Printf("successful: %d\n", len(results))
	fmt.Printf("errors: %d\n", len(errors))

	for _, e := range errors {
		fmt.Println("error:", e)
	}
}
