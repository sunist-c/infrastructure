package task

type Status int8

const (
	StatusCreated   Status = 0
	StatusRunning   Status = 1
	StatusRetrying  Status = 2
	StatusCancelled Status = 3
	StatusFailed    Status = 4
	StatusSucceeded Status = 5
	StatusLooped    Status = 6
)
