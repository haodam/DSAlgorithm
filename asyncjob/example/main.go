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

	job2 := asyncjob.NewJob(func(ctx context.Context) error {
		time.Sleep(time.Second * 2)
		log.Println("I am job 2")

		return nil
		//return errors.New("something went wrong at job 2")
	})

	job3 := asyncjob.NewJob(func(ctx context.Context) error {
		time.Sleep(time.Second * 3)
		log.Println("I am job 3")

		return nil
		//return errors.New("something went wrong at job 3")
	})

	group := asyncjob.NewGroup(false, job1, job2, job3)
	if err := group.Run(context.Background()); err != nil {
		log.Println(err)
	}

	//if err := job1.Execute(context.Background()); err != nil {
	//	log.Println(job1.State(), err)
	//}
	//
	//for {
	//	if err := job1.Retry(context.Background()); err != nil {
	//		log.Println(err)
	//	}
	//	if job1.State() == asyncjob.JobStateRetryFailed || job1.State() == asyncjob.JobStateCompleted {
	//		log.Println(job1.State())
	//		break
	//	}
	//}
}
