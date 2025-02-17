package logger

type entry struct {
	File     string         `json:"file"`
	Level    string         `json:"level"`
	Service  string         `json:"service"`
	TraceID  string         `json:"trace_id"`
	Message  string         `json:"message"`
	CallTime string         `json:"call_time"`
	Data     any            `json:"data"`
	Extra    map[string]any `json:"extra,omitempty"`
}
