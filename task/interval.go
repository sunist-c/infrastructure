package task

import "time"

func StaticInterval(interval time.Duration, maxRetries int) RetryInterval {
	return func(retries int) (bool, time.Duration) {
		if retries > maxRetries {
			return false, interval
		}

		return true, interval
	}
}

func BackoffInterval(baseInterval time.Duration, maxRetries int) RetryInterval {
	return func(retries int) (needRetry bool, interval time.Duration) {
		if retries > maxRetries {
			return false, baseInterval
		}

		return true, baseInterval * time.Duration(retries)
	}
}

func LoopInterval(baseInterval time.Duration) RetryInterval {
	return func(retries int) (needRetry bool, interval time.Duration) {
		return true, baseInterval
	}
}
