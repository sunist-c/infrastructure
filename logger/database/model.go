package adlog

type LogContent struct {
	Trace        string `json:"trace"`
	Level        string `json:"level"`
	Service      string `json:"service"`
	Instance     string `json:"instance"`
	Message      string `json:"message"`
	SQL          string `json:"sql,omitempty"`
	RowsAffected int64  `json:"rows_affected,omitempty"`
	ExecError    string `json:"exec_error,omitempty"`
	ExecTime     string `json:"exec_time,omitempty"`
}
