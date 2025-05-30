package asyncjob

import (
	"context"
	"log"
	"sync"
)

type group struct {
	isConcurrent bool
	jobs         []Job
	wg           *sync.WaitGroup
}

func NewGroup(isConcurrent bool, jobs ...Job) *group {
	g := &group{
		isConcurrent: isConcurrent,
		jobs:         jobs,
		wg:           new(sync.WaitGroup),
	}
	return g
}

func (g *group) Run(ctx context.Context) error {

	g.wg.Add(len(g.jobs))

	errChan := make(chan error, len(g.jobs))
	for i, _ := range g.jobs {
		if g.isConcurrent {
			go func(aj Job) {
				errChan <- g.RunJob(ctx, aj)
				g.wg.Done()
			}(g.jobs[i])
			continue
		}

		job := g.jobs[i]
		errChan <- g.RunJob(ctx, job)
		g.wg.Done()
	}

	var err error
	for i := 0; i < len(g.jobs); i++ {
		if v := <-errChan; v != nil {
			err = v
		}
	}
	g.wg.Wait()
	return err
}

func (g *group) RunJob(ctx context.Context, j Job) error {
	if err := j.Execute(ctx); err != nil {
		for {
			log.Println(err)
			if j.State() == JobStateRetryFailed {
				return err
			}
			if j.Retry(ctx) == nil {
				return nil
			}
		}
	}
	return nil
}
