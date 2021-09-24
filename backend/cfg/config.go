package cfg

// 应用配置
const (
	AppName = "cantor"
	Version = "v0.0.6"
)

// 图片配置
var (
	AllowFileExts       = []string{".png", ".gif", ".jpg", ".jpeg"}
	MaxFileSize   int64 = 2 * 1024 * 1024
)

// Git 配置
var (
	GitApiUrl   = "https://api.github.com/repos/%s/%s/contents/%s"
	GitTagUrl   = "https://api.github.com/repos/evercyan/cantor/tags"
	GitFileUrl  = "https://cdn.jsdelivr.net/gh/%s/%s/%s"
	GitDBFile   = "resource/cantor.db"
	GitFilePath = "resource/%s/%s%s"
	GitMessage  = "upload by cantor"
)

// 文件配置
var (
	CfgFile = "%s/config.json"
	LogFile = "%s/app.log"
	DBFile  = "%s/cantor.db"
)

// Resp ...
type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
