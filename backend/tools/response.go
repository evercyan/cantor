package tools

import (
	"github.com/evercyan/cantor/backend/configs"
)

// Success ...
func Success(data interface{}) *configs.Resp {
	return &configs.Resp{
		Code: 0,
		Data: data,
	}
}

// Fail ...
func Fail(message string) *configs.Resp {
	return &configs.Resp{
		Code: -1,
		Msg:  message,
	}
}
