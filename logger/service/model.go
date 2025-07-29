package aslog

type LogContent struct {
	Trace       string            `json:"trace"`
	Level       string            `json:"level"`
	LogPoint    string            `json:"log_point"`
	Service     string            `json:"service"`
	Instance    string            `json:"instance"`
	Timestamp   string            `json:"ts"`
	TimeCost    string            `json:"tc"`
	Message     string            `json:"message"`
	Content     any               `json:"content,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	RequestInfo *RequestContent   `json:"request,omitempty"`
}

type RequestContent struct {
	Method   string `json:"method"`
	Path     string `json:"path"`
	Host     string `json:"host"`
	ClientIP string `json:"client_ip"`
}
