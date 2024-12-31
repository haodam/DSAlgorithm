package main

import (
	"context"
	"errors"
	"github.com/haodam/DSAlgorithm/asyncjob"
	"log"
	"time"
)

func main() {
	job1 := asyncjob.NewJob(func(ctx context.Context) error {
		time.Sleep(time.Second)
		log.Println("I am job 1")

		//return nil
		return errors.New("something went wrong at job 1")
	})

	if err := job1.Execute(context.Background()); err != nil {
		log.Println(job1.State(), err)
	}

	for {
		if err := job1.Retry(context.Background()); err != nil {
			log.Println(err)
		}
		if job1.State() == asyncjob.JobStateRetryFailed || job1.State() == asyncjob.JobStateCompleted {
			log.Println(job1.State())
			break
		}
	}
}
