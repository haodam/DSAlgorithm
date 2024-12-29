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

const (
	defaultMaxTimeout = 10 * time.Second
)

var (
	defaultRetryTime = []time.Duration{time.Second, time.Second * 5, time.Second * 10}
)

type JobHandler func(ctx context.Context, job *Job) error

type JobState int

const (
	JobStateInt JobState = iota
	JobStateRunning
	JobStateFailed
	JobStateCompleted
	JobStateRetryFailed
)

type Job interface {
	Execute(ctx context.Context) error
	Retry(ctx context.Context) error
	State() JobState
	SetRetryDuration(times []time.Duration)
}
