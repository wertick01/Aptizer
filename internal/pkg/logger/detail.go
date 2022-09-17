package logger

import (
	"fmt"
	"runtime/debug"
)

type Detail struct {
	Request      string         `json:"request,omitempty"`
	Address      string         `json:"address,omitempty"`
	Backtrace    string         `json:"backtrace,omitempty"`
	Response     string         `json:"response,omitempty"`
	ResponseCode int            `json:"response_code,omitempty"`
	Fields       map[string]any `json:"fields,omitempty"`
}

func GetBacktrace(err error) string {
	stack := string(debug.Stack())

	if err != nil {
		stack = fmt.Sprintf("%s\n%s", err.Error(), stack)
	}

	return stack
}
