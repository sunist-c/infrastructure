package trace

import (
	"context"
)

var (
	instance string
	service  string

	background = Trace(context.Background())
)

type Type string

func (t Type) String() string {
	return string(t)
}

const (
	BasicType   Type = "at:basic_info"
	RequestType Type = "at:request_info"
)
