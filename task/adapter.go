package task

import (
	"time"
)

type RetryInterval func(retries int) (needRetry bool, interval time.Duration)

type Manager interface {
	Execute(fn func() bool, retryInterval RetryInterval) (cancel func())
	ExecuteLoop(fn func(), loopInterval time.Duration) (cancel func())
}
