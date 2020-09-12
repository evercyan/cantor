package backend

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/d-tsuji/clipboard"
	"github.com/evercyan/letitgo/crypto"
	ljson "github.com/evercyan/letitgo/json"
	"github.com/evercyan/letitgo/util"
	"github.com/wailsapp/wails"
)

var (
	allowFileExts       = []string{".png", ".gif", ".jpg", ".jpeg"}
	maxFileSize   int64 = 2 * 1024 * 1024
	gitDBFile           = "resource/database.json"
	gitFilePath         = "resource/%s/%s%s"
)

/**
 * ********************************
 */

// App ...
type App struct {
	RT         *wails.Runtime
	ConfigFile string
	Git        Git
}

// WailsInit ...
func (a *App) WailsInit(runtime *wails.Runtime) error {
	Log().Info("WailsInit")
	a.RT = runtime
	a.initConfig()
	return nil
}

// WailsShutdown ...
func (a *App) WailsShutdown() {
	Log().Info("WailsShutdown")
	return
}

// 配置初始化
func (a *App) initConfig() {
	user, err := user.Current()
	if err != nil {
		panic("获取用户目录失败: " + err.Error())
	}

	// 配置目录
	cantorPath := user.HomeDir + "/.cantor"
	if !util.IsExist(cantorPath) {
		os.Mkdir(cantorPath, os.ModePerm)
	}

	// 配置文件
	a.ConfigFile = cantorPath + "/config.json"
	Log().Info("initConfig configFile ", a.ConfigFile)
	configContent := util.ReadFile(a.ConfigFile)
	Log().Info("GetConfig configContent ", configContent)
	if configContent == "" {
		return
	}
	json.Unmarshal([]byte(configContent), &a.Git)
}

// Resp 返回封装
func (a *App) Resp(code int, data interface{}) string {
	return crypto.JsonEncode(map[string]interface{}{
		"code": code,
		"data": data,
	})
}

/**
 * ********************************
 */

func (a *App) getUploadList() []map[string]string {
	resp, _ := a.Git.Get(gitDBFile)
	list := []map[string]string{}
	if resp != "" {
		content := crypto.Base64Decode(ljson.Json(resp).Key("content").ToString())
		json.Unmarshal([]byte(content), &list)
	}
	return list
}

/**
 * ********************************
 */

// GetConfig 获取 git 配置
func (a *App) GetConfig(param string) string {
	return a.Resp(0, a.Git)
}

// SetConfig 更新 git 配置
func (a *App) SetConfig(param string) string {
	Log().Info("SetConfig param ", param)
	if err := json.Unmarshal([]byte(param), &a.Git); err != nil {
		return a.Resp(-1, err.Error())
	}
	if err := util.WriteFile(a.ConfigFile, param); err != nil {
		return a.Resp(-1, err.Error())
	}
	return a.Resp(0, "设置成功")
}

// GetUploadList 获取上传文件列表
func (a *App) GetUploadList(param string) string {
	return a.Resp(0, a.getUploadList())
}

// UploadFile 上传文件
func (a *App) UploadFile(param string) string {
	selectFile := a.RT.Dialog.SelectFile()
	Log().Info("UploadFile selectFile ", selectFile)
	if selectFile == "" {
		return a.Resp(-1, "请选择文件")
	}
	if a.Git.Repo == "" {
		return a.Resp(-1, "请设置配置")
	}

	// 文件格式校验
	fileExt := strings.ToLower(path.Ext(selectFile))
	if !util.InArray(fileExt, allowFileExts) {
		return a.Resp(-1, "仅支持以下文件格式: "+strings.Join(allowFileExts, ", "))
	}

	// 文件大小校验
	fileSize := util.GetSize(selectFile)
	if fileSize > maxFileSize {
		return a.Resp(-1, "最大支持 2M 的文件")
	}

	// 文件内容
	fileContent := util.ReadFile(selectFile)
	// 文件路径名称
	fileMd5 := util.Md5(fileContent)
	filePath := fmt.Sprintf(gitFilePath, fileMd5[0:2], fileMd5, fileExt)
	// 请求上传文件
	err := a.Git.Update(filePath, fileContent)
	if err != nil {
		return a.Resp(-1, err.Error())
	}

	// 更新数据文件
	fileUrl := a.Git.Url(filePath)
	fileName := path.Base(selectFile)
	fileInfo := map[string]string{
		"file_name": fileName,
		"file_md5":  fileMd5,
		"file_size": util.GetSizeText(fileSize),
		"file_path": filePath,
		"file_url":  fileUrl,
		"create_at": time.Now().Format("2006-01-02 15:04:05"),
	}
	Log().Info("UploadFile fileInfo ", crypto.JsonEncode(fileInfo))
	list := a.getUploadList()
	fileMd5List := []string{}
	for i := 0; i < len(list); i++ {
		fileMd5List = append(fileMd5List, list[i]["file_md5"])
	}
	if !util.InArray(fileMd5, fileMd5List) {
		list = append([]map[string]string{fileInfo}, list...)
		updateErr := a.Git.Update(gitDBFile, crypto.JsonEncode(list))
		if updateErr != nil {
			Log().Info("UploadFile updateErr ", updateErr.Error())
		}
	}

	return a.Resp(0, fileInfo)
}

// DeleteFile 删除文件
func (a *App) DeleteFile(param string) string {
	list := a.getUploadList()
	for i := 0; i < len(list); i++ {
		if list[i]["file_path"] == param {
			if i == len(list)-1 {
				list = list[:i]
			} else {
				list = append(list[:i], list[i+1:]...)
			}
		}
	}
	// 更新数据文件
	updateErr := a.Git.Update(gitDBFile, crypto.JsonEncode(list))
	if updateErr != nil {
		return a.Resp(-1, updateErr.Error())
	}
	// 删除文件
	deleteErr := a.Git.Delete(param)
	if deleteErr != nil {
		Log().Info("UploadFile deleteErr ", deleteErr.Error())
	}
	return a.Resp(0, "")
}

// CopyFileUrl 复制链接
func (a *App) CopyFileUrl(fileUrl string) string {
	Log().Info("CopyFileUrl fileUrl ", fileUrl)
	err := clipboard.Set(fileUrl)
	if err != nil {
		Log().Info("CopyFileUrl err ", err.Error())
		return a.Resp(-1, err.Error())
	}
	return a.Resp(0, "已复制到粘贴板")
}
