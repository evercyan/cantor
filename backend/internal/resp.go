package internal

// Resp ...
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Success ...
func Success(data interface{}) *Response {
	return &Response{
		Code: 0,
		Data: data,
	}
}

// Fail ...
func Fail(message string) *Response {
	return &Response{
		Code: -1,
		Msg:  message,
	}
}
