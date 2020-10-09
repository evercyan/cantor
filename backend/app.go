package backend

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/d-tsuji/clipboard"
	"github.com/evercyan/cantor/backend/configs"
	"github.com/evercyan/cantor/backend/internal/git"
	"github.com/evercyan/cantor/backend/tools"
	"github.com/evercyan/letitgo/crypto"
	"github.com/evercyan/letitgo/util"
	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails"
)

// App ...
type App struct {
	RT     *wails.Runtime
	Log    *logrus.Logger
	Git    git.Git
	Config string
}

// WailsInit ...
func (a *App) WailsInit(runtime *wails.Runtime) error {
	a.RT = runtime
	a.Log = tools.NewLog()
	a.Log.Info("WailsInit")
	a.InitConfig()
	return nil
}

// WailsShutdown ...
func (a *App) WailsShutdown() {
	a.Log.Info("WailsShutdown")
	return
}

// InitConfig ...
func (a *App) InitConfig() {
	a.Config = tools.GetConfigPath() + "/config.json"
	a.Log.Info("InitConfig config: ", a.Config)
	content := util.ReadFile(a.Config)
	a.Log.Info("InitConfig content: ", content)
	if content == "" {
		return
	}
	json.Unmarshal([]byte(content), &a.Git)
}

// GetConfig 获取 git 配置和版本信息
func (a *App) GetConfig() *configs.Resp {
	resp := map[string]interface{}{
		"config": a.Git,
		"version": map[string]interface{}{
			"current": configs.Version,
			"last":    a.Git.LastVersion(),
		},
	}
	a.Log.Info("GetConfig content: ", crypto.JsonEncode(resp))
	return tools.Success(resp)
}

// SetConfig 更新 git 配置
func (a *App) SetConfig(content string) *configs.Resp {
	a.Log.Info("SetConfig content: ", content)
	if err := json.Unmarshal([]byte(content), &a.Git); err != nil {
		return tools.Fail(err.Error())
	}
	if err := util.WriteFile(a.Config, content); err != nil {
		return tools.Fail(err.Error())
	}
	return tools.Success("操作成功")
}

// GetUploadList 获取上传文件列表
func (a *App) GetUploadList() *configs.Resp {
	return tools.Success(a.Git.UploadFileList())
}

// UploadFile 上传文件
func (a *App) UploadFile() *configs.Resp {
	selectFile := a.RT.Dialog.SelectFile()
	a.Log.Info("UploadFile selectFile ", selectFile)
	if selectFile == "" {
		return tools.Fail("请选择图片文件")
	}
	if a.Git.Repo == "" {
		return tools.Fail("请设置 Git 配置")
	}

	// 文件格式校验
	fileExt := strings.ToLower(path.Ext(selectFile))
	if !util.InArray(fileExt, configs.AllowFileExts) {
		return tools.Fail("仅支持以下格式: " + strings.Join(configs.AllowFileExts, ", "))
	}

	// 文件大小校验
	fileSize := util.GetSize(selectFile)
	if fileSize > configs.MaxFileSize {
		return tools.Fail("最大支持 2M 的文件")
	}

	// 文件内容
	fileContent := util.ReadFile(selectFile)
	// 文件路径名称
	fileMd5 := util.Md5(fileContent)
	filePath := fmt.Sprintf(configs.GitFilePath, fileMd5[0:2], fileMd5, fileExt)
	// 请求上传文件
	err := a.Git.Update(filePath, fileContent)
	if err != nil {
		return tools.Fail(err.Error())
	}

	// 更新数据文件
	fileInfo := map[string]string{
		"file_name": path.Base(selectFile),
		"file_md5":  fileMd5,
		"file_size": util.GetSizeText(fileSize),
		"file_path": filePath,
		"file_url":  a.Git.Url(filePath),
		"create_at": time.Now().Format("2006-01-02 15:04:05"),
	}
	a.Log.Info("UploadFile fileInfo: ", crypto.JsonEncode(fileInfo))
	list := a.Git.UploadFileList()
	list = append([]map[string]string{fileInfo}, list...)
	updateErr := a.Git.Update(configs.GitDBFile, crypto.JsonEncode(list))
	if updateErr != nil {
		a.Log.Info("UploadFile updateErr: ", updateErr.Error())
	}
	return tools.Success("操作成功")
}

// DeleteFile 删除文件
func (a *App) DeleteFile(filePath string) *configs.Resp {
	list := a.Git.UploadFileList()
	for i := 0; i < len(list); i++ {
		if list[i]["file_path"] == filePath {
			if i == len(list)-1 {
				list = list[:i]
			} else {
				list = append(list[:i], list[i+1:]...)
			}
		}
	}
	// 更新数据文件
	updateErr := a.Git.Update(configs.GitDBFile, crypto.JsonEncode(list))
	if updateErr != nil {
		return tools.Fail(updateErr.Error())
	}
	// 删除文件
	deleteErr := a.Git.Delete(filePath)
	if deleteErr != nil {
		a.Log.Info("UploadFile deleteErr: ", deleteErr.Error())
	}
	return tools.Success("操作成功")
}

// CopyFileUrl 复制链接到粘贴板
func (a *App) CopyFileUrl(fileUrl string) *configs.Resp {
	a.Log.Info("CopyFileUrl fileUrl: ", fileUrl)
	err := clipboard.Set(fileUrl)
	if err != nil {
		return tools.Fail(err.Error())
	}
	return tools.Success("已复制到粘贴板")
}
