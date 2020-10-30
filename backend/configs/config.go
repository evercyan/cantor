package configs

const (
	// AppName 应用名称
	AppName = "cantor"
	// Version 应用版本
	Version = "v0.0.5"
)

var (
	// AllowFileExts 图片格式限制
	AllowFileExts = []string{".png", ".gif", ".jpg", ".jpeg"}
	// MaxFileSize 图片大小限制
	MaxFileSize int64 = 2 * 1024 * 1024
)

var (
	// GitApiUrl ...
	GitApiUrl = "https://api.github.com/repos/%s/%s/contents/%s?access_token=%s"
	// GitTagUrl ...
	GitTagUrl = "https://api.github.com/repos/%s/%s/tags"
	// GitFileUrl ...
	GitFileUrl = "https://raw.githubusercontent.com/%s/%s/master/%s"
	// GitDBFile ...
	GitDBFile = "resource/database.json"
	// GitFilePath ...
	GitFilePath = "resource/%s/%s%s"
	// GitMessage ...
	GitMessage = "auto deploy"
)

// Resp ...
type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
