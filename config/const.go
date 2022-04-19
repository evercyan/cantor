package config

// 应用
const (
	App     = "Cantor"
	Version = "v0.1.0"
)

const (
	Title          = App + " " + Version
	Description    = "一个简单好用的图床应用"
	VersionNewMsg  = "当前已经是最新版本!"
	VersionOldMsg  = "最新版本: %s, 是否立即更新?"
	BtnConfirmText = "确定"
	BtnCancelText  = "取消"
)

// 窗口尺寸
const (
	Width  = 1024
	Height = 768
)

// 图片配置
var (
	AllowFileExts       = []string{".png", ".gif", ".jpg", ".jpeg"}
	MaxFileSize   int64 = 4 * 1024 * 1024
	MaxFileCount        = 10
)

// 文件配置
var (
	CfgFile = "%s/config.yaml"
	LogFile = "%s/app.log"
	DBFile  = "%s/cantor.db"
)
