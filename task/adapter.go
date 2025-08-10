package task

import (
	"time"

	"github.com/alioth-center/infrastructure/grace"
)

type RetryInterval func(retries int) (needRetry bool, interval time.Duration)

type Trigger interface {
	Trigger() (Status, func())
}

type Manager interface {
	Execute(fn func() bool, retryInterval RetryInterval) (cancel func())
	ExecuteLoop(fn func(), loopInterval time.Duration) (cancel func())
	ExecuteWithTriggers(fn func() bool, retryInterval RetryInterval, triggers ...Trigger) (cancel func())
	ExecuteLoopWithTriggers(fn func(), loopInterval time.Duration, triggers ...Trigger) (cancel func())
}

type GracefulManager interface {
	grace.Graceful
	Manager
}
