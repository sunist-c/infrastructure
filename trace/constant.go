package trace

type Type string

func (t Type) String() string {
	return string(t)
}

const (
	BasicType   Type = "at:basic_info"
	RequestType Type = "at:request_info"
)
