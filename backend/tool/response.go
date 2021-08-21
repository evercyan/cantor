package tool

import (
	"github.com/evercyan/cantor/backend/cfg"
)

// Success ...
func Success(data interface{}) *cfg.Resp {
	return &cfg.Resp{
		Code: 0,
		Data: data,
	}
}

// Fail ...
func Fail(message string) *cfg.Resp {
	return &cfg.Resp{
		Code: -1,
		Msg:  message,
	}
}
