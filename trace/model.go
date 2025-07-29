package trace

import "time"

type Basic struct {
	TraceID  string
	TracedAt time.Time
	Instance string
}

type Request struct {
	Method   string
	Path     string
	Host     string
	ClientIP string
}
