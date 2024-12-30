package asyncjob

import (
	"context"
	"time"
)

// Job Requirement
// 1.   Job can do something (handler)
// 2.   Job can retry
// 2.1  Config retry time and duration
// 3.   Should be stateful
// 4.   We should have job manager to manage jobs (***)

type Job interface {
	Execute(ctx context.Context) error
	Retry(ctx context.Context) error
	State() JobState
	SetRetryDuration(times []time.Duration)
}

const (
	defaultMaxTimeout = 10 * time.Second
)

var (
	defaultRetryTime = []time.Duration{time.Second, time.Second * 5, time.Second * 10}
)

type JobHandler func(ctx context.Context) error

type JobState int

const (
	JobStateInt JobState = iota
	JobStateRunning
	JobStateFailed
	JobStateTimeout
	JobStateCompleted
	JobStateRetryFailed
)

func (js JobState) String() string {
	return []string{"Init", "Running", "Failed", "Timeout", "Completed"}[js]
}

type JobConfig struct {
	MaxTimeout time.Duration
	Retries    []time.Duration
}

type job struct {
	config     JobConfig
	handler    JobHandler
	state      JobState
	retryIndex int
	stopChan   chan bool
}

func NewJob(handler JobHandler) *job {
	j := job{
		config: JobConfig{
			MaxTimeout: defaultMaxTimeout,
			Retries:    defaultRetryTime,
		},
		handler:    handler,
		state:      JobStateInt,
		retryIndex: -1, // chua retry lan nao het
		stopChan:   make(chan bool),
	}
	return &j
}

func (j *job) Execute(ctx context.Context) error {

	j.state = JobStateInt
	var err error
	err = j.handler(ctx)

	if err != nil {
		j.state = JobStateFailed
		return err
	}

	j.state = JobStateCompleted
	return nil
}

func (j *job) Retry(ctx context.Context) error {
	return nil
}

func (j *job) State() JobState {
	return 0
}

func (j *job) SetRetryDuration(times []time.Duration) {

}
