package backend

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/d-tsuji/clipboard"
	"github.com/evercyan/cantor/backend/cfg"
	"github.com/evercyan/cantor/backend/pkg"
	"github.com/evercyan/cantor/backend/tool"
	"github.com/evercyan/letitgo/crypto"
	"github.com/evercyan/letitgo/file"
	lj "github.com/evercyan/letitgo/json"
	"github.com/evercyan/letitgo/util"
	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails"
	"gorm.io/gorm"
)

// App ...
type App struct {
	RT      *wails.Runtime
	Log     *logrus.Logger
	Git     *pkg.Git
	DB      *gorm.DB
	CfgFile string
	DBFile  string
}

// ----------------------------------------------------------------

// WailsInit ...
func (a *App) WailsInit(runtime *wails.Runtime) error {
	a.RT = runtime
	// 日志
	a.Log = tool.NewLogger()
	a.Log.Info("WailsInit")
	// 配置
	cfgPath := tool.GetConfigPath()
	a.CfgFile = fmt.Sprintf(cfg.CfgFile, cfgPath)
	a.Git = &pkg.Git{}
	if err := json.Unmarshal([]byte(file.Read(a.CfgFile)), a.Git); err != nil {
		a.Log.Errorf("WailsInit CfgFile err: %v", err)
	}
	a.Log.Infof("WailsInit CfgFile: %+v", a.Git)
	// db
	a.DBFile = fmt.Sprintf(cfg.DBFile, cfgPath)
	isExist := file.IsExist(a.DBFile)
	a.DB = pkg.NewDB(a.DBFile)
	if !isExist && a.Git != nil {
		// 如果之前文件不存在, 需要做 db 迁移
		a.migrate()
	}
	a.Log.Info("WailsInit success")
	return nil
}

// migrate ...
func (a *App) migrate() {
	// 获取 v0.0.5 及之前存储上传的文件纪录, 迁移至 sqlite db
	resp, _ := a.Git.Get("resource/database.json")
	if resp == "" {
		a.Log.Error("WailsInit migrate err: empty database.json")
		return
	}
	list := make([]pkg.File, 0)
	content := crypto.Base64Decode(lj.Json(resp).Key("content").ToString())
	if err := json.Unmarshal([]byte(content), &list); err != nil {
		a.Log.Errorf("WailsInit migrate err: %v", err)
		return
	}
	a.Log.Infof("WailsInit migrate total count: %d", len(list))
	// 迁移到本地 sqlite db
	success := 0
	for i := len(list) - 1; i >= 0; i-- {
		res := a.DB.Create(&list[i])
		if res.Error != nil {
			a.Log.Errorf("WailsInit migrate file: %s, err: %v", list[i].Name, res.Error)
		} else {
			success++
		}
	}
	a.Log.Infof("WailsInit migrate success count: %d", success)
	return
}

// WailsShutdown ...
func (a *App) WailsShutdown() {
	a.Log.Info("WailsShutdown")
	return
}

// ----------------------------------------------------------------

// GetConfig 获取配置信息
func (a *App) GetConfig() *cfg.Resp {
	resp := map[string]interface{}{
		"config": a.Git,
		"version": map[string]interface{}{
			"current": cfg.Version,
			"last":    a.Git.GetLastVersion(),
		},
	}
	a.Log.Infof("GetConfig resp: %s", crypto.JsonEncode(resp))
	return tool.Success(resp)
}

// SetConfig 更新配置信息
func (a *App) SetConfig(content string) *cfg.Resp {
	a.Log.Infof("SetConfig content: %v", content)
	if err := json.Unmarshal([]byte(content), a.Git); err != nil {
		a.Log.Errorf("SetConfig Git err: %v", err)
		return tool.Fail(err.Error())
	}
	if err := file.Write(a.CfgFile, content); err != nil {
		a.Log.Errorf("SetConfig Write err: %v", err)
		return tool.Fail(err.Error())
	}
	return tool.Success("操作成功")
}

// ----------------------------------------------------------------

// GetList 获取文件列表
func (a *App) GetList() *cfg.Resp {
	fileList := make([]pkg.File, 0)
	a.DB.Order("id DESC").Find(&fileList)
	a.Log.Infof("GetList count: %v", len(fileList))
	if len(fileList) > 0 {
		a.Log.Infof("GetList[0]: %s", crypto.JsonEncode(fileList[0]))
	}
	list := make([]map[string]string, 0)
	for _, fileInfo := range fileList {
		list = append(list, map[string]string{
			"file_name": fileInfo.Name,
			"file_md5":  fileInfo.Md5,
			"file_size": fileInfo.Size,
			"file_path": fileInfo.Path,
			"file_url":  a.Git.GetFileUrl(fileInfo.Path),
			"create_at": fileInfo.CreateAt,
		})
	}
	return tool.Success(list)
}

// UploadFile 上传文件
func (a *App) UploadFile() *cfg.Resp {
	selectFile := a.RT.Dialog.SelectFile()
	a.Log.Infof("UploadFile selectFile: %v", selectFile)
	if selectFile == "" {
		return tool.Fail("请选择图片文件")
	}
	if a.Git.Repo == "" {
		return tool.Fail("请设置 Git 配置")
	}
	// 文件格式校验
	fileExt := strings.ToLower(path.Ext(selectFile))
	if !util.InArray(fileExt, cfg.AllowFileExts) {
		return tool.Fail("仅支持以下格式: " + strings.Join(cfg.AllowFileExts, ", "))
	}
	// 文件大小校验
	fileSize := file.Size(selectFile)
	a.Log.Infof("UploadFile fileSize: %v", fileSize)
	if fileSize > cfg.MaxFileSize {
		return tool.Fail("最大支持 2M 的文件")
	}
	fileContent := file.Read(selectFile)
	fileMd5 := util.Md5(fileContent)
	filePath := fmt.Sprintf(cfg.GitFilePath, fileMd5[0:2], fileMd5, fileExt)
	// 请求上传文件
	err := a.Git.Update(filePath, fileContent)
	if err != nil {
		return tool.Fail(err.Error())
	}
	// 新增数据
	fileInfo := &pkg.File{
		Name:     path.Base(selectFile),
		Md5:      fileMd5,
		Size:     file.SizeText(fileSize),
		Path:     filePath,
		CreateAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	a.Log.Infof("UploadFile fileInfo: %+v", fileInfo)
	res := a.DB.Create(fileInfo)
	if res.Error != nil {
		a.Log.Errorf("UploadFile DB err: %v", res.Error)
		return tool.Fail(res.Error.Error())
	}
	go func() {
		// 上传本地 db 文件
		err := a.Git.Update(cfg.GitDBFile, file.Read(a.DBFile))
		if err != nil {
			a.Log.Errorf("UploadFile Update DBFile err: %v", err)
		} else {
			a.Log.Info("UploadFile Update DBFile success")
		}
	}()
	a.Log.Info("UploadFile success")
	return tool.Success("操作成功")
}

// DeleteFile 删除文件
func (a *App) DeleteFile(filePath string) *cfg.Resp {
	a.Log.Infof("DeleteFile filePath: %v", filePath)
	// 删除文件
	err := a.Git.Delete(filePath)
	if err != nil {
		a.Log.Errorf("DeleteFile Git err: %v", err)
	}
	// 更新 DB
	res := a.DB.Where("file_path", filePath).Delete(&pkg.File{})
	if res.Error != nil {
		a.Log.Errorf("DeleteFile DB err: %v", res.Error)
		return tool.Fail(res.Error.Error())
	}
	a.Log.Info("DeleteFile success")
	return tool.Success("操作成功")
}

// UpdateFileName 更新文件名称
func (a *App) UpdateFileName(filePath string, fileName string) *cfg.Resp {
	a.Log.Infof("UpdateFileName filePath: %v; fileName: %v", filePath, fileName)
	if fileName == "" {
		return tool.Fail("文件名称不能为空")
	}
	res := a.DB.Model(&pkg.File{}).Where("file_path", filePath).Update("file_name", fileName)
	if res.Error != nil {
		a.Log.Errorf("UpdateFileName DB err: %v", res.Error)
		return tool.Fail(res.Error.Error())
	}
	a.Log.Info("UpdateFileName success")
	return tool.Success("操作成功")
}

// ----------------------------------------------------------------

// CopyFileUrl 复制链接到粘贴板
func (a *App) CopyFileUrl(fileUrl string) *cfg.Resp {
	a.Log.Info("CopyFileUrl fileUrl: ", fileUrl)
	err := clipboard.Set(fileUrl)
	if err != nil {
		return tool.Fail(err.Error())
	}
	a.Log.Info("CopyFileUrl success")
	return tool.Success("已复制到粘贴板")
}
