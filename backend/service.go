package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/atotto/clipboard"
	"github.com/evercyan/brick/xcrypto"
	"github.com/evercyan/brick/xencoding"
	"github.com/evercyan/brick/xfile"
	"github.com/evercyan/brick/xtime"
	"github.com/evercyan/cantor/backend/internal"
	"github.com/evercyan/cantor/config"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v2"
)

// GetConfig 获取配置信息
func (t *App) GetConfig() *internal.Response {
	t.Log.Infof("GetConfig resp: %s", xencoding.JSONEncode(t.Git))
	return internal.Success(t.Git)
}

// ----------------------------------------------------------------

// SetConfig 更新配置信息
func (t *App) SetConfig(content string) *internal.Response {
	t.Log.Infof("SetConfig content: %v", content)
	git := &internal.Git{}
	if err := json.Unmarshal([]byte(content), git); err != nil {
		t.Log.Errorf("SetConfig Unmarshal err: %v", err)
		return internal.Fail(err.Error())
	}
	if err := git.Update(config.GitMarkFile, xtime.Format(time.Now(), "ymdhis")); err != nil {
		t.Log.Errorf("SetConfig GitHub err: %v", err)
		return internal.Fail("无效的 Git 配置")
	}
	t.Git = git
	b, _ := yaml.Marshal(t.Git)
	if err := xfile.Write(t.CfgFile, string(b)); err != nil {
		t.Log.Errorf("SetConfig Write err: %v", err)
		return internal.Fail(err.Error())
	}
	// 数据库
	t.database()
	return internal.Success("操作成功")
}

// ----------------------------------------------------------------

// GetList 获取文件列表
func (t *App) GetList() *internal.Response {
	t.Log.Info("GetList begin")
	if t.DB == nil {
		// 特殊处理 OnStartup 时如果存在数据库迁移, 可能出现未生成 t.DB 的情况
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
		defer cancelFunc()
		for t.DB == nil {
			time.Sleep(time.Second * 1)
			if ctx.Err() != nil {
				return internal.Fail("数据库迁移失败, 请稍候重试")
			}
		}
	}
	fileList := make([]internal.File, 0)
	t.DB.Order("create_at DESC").Find(&fileList)
	t.Log.Infof("GetList count: %v", len(fileList))
	if len(fileList) > 0 {
		t.Log.Infof("GetList first: %s", xencoding.JSONEncode(fileList[0]))
	}
	list := make([]map[string]string, 0)
	for _, fileInfo := range fileList {
		list = append(list, map[string]string{
			"file_name": fileInfo.Name,
			"file_md5":  fileInfo.Md5,
			"file_size": fileInfo.Size,
			"file_path": fileInfo.Path,
			"file_url":  t.Git.GetFileUrl(fileInfo.Path),
			"create_at": fileInfo.CreateAt,
		})
	}
	return internal.Success(list)
}

// ----------------------------------------------------------------

// BatchUploadFile ...
func (t *App) BatchUploadFile() *internal.Response {
	if t.Git.Repo == "" {
		return internal.Fail("请先更新配置")
	}
	files, err := runtime.OpenMultipleFilesDialog(t.Ctx, runtime.OpenDialogOptions{
		Title: "选择图片",
		Filters: []runtime.FileFilter{{
			DisplayName: "Images (*.png;*.jpg;*.jpeg;*.gif)",
			Pattern:     "*.png;*.jpg;*.jpeg;*.gif",
		}},
	})
	if err != nil {
		return internal.Fail(err.Error())
	}
	if len(files) == 0 {
		return internal.Fail("请选择至少一张图片")
	}
	if len(files) > config.MaxFileCount {
		return internal.Fail(fmt.Sprintf("最多可选择 %d 张图片", config.MaxFileCount))
	}
	for _, file := range files {
		err := t.CheckFile(file)
		if err != nil {
			return internal.Fail(fmt.Sprintf("%s: %s", path.Base(file), err.Error()))
		}
	}

	var wg sync.WaitGroup
	var mx sync.Mutex
	count := 0
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			_, err := t.Upload(file)
			if err != nil {
				t.Log.Errorf("BatchUploadFile file: %s , err: %s", file, err.Error())
				return
			}
			mx.Lock()
			count++
			mx.Unlock()
		}(file)
	}
	wg.Wait()

	if count == 0 {
		return internal.Fail("上传图片失败")
	}
	go t.SyncDatabase()
	return internal.Success(fmt.Sprintf("上传图片 %d 张, 成功 %d 张", len(files), count))
}

// ----------------------------------------------------------------

// Upload ...
func (t *App) Upload(filepath string, clis ...bool) (string, error) {
	t.Log.Infof("Upload filepath: %v", filepath)
	fileExt := strings.ToLower(path.Ext(filepath))
	fileSize := xfile.Size(filepath)
	fileContent := xfile.Read(filepath)
	fileMd5 := xcrypto.Md5(fileContent)
	fileGitPath := fmt.Sprintf(config.GitFilePath, fileMd5[0:2], fileMd5, fileExt)
	// 请求上传文件
	err := t.Git.Update(fileGitPath, fileContent)
	if err != nil {
		t.Log.Infof("Upload Git err: %v", err)
		return "", err
	}
	// 新增/更新数据
	fileInfo := &internal.File{}
	t.DB.Where("file_md5", fileMd5).First(fileInfo)
	fileInfo.Name = path.Base(filepath)
	fileInfo.Md5 = fileMd5
	fileInfo.Size = xfile.SizeText(fileSize)
	fileInfo.Path = fileGitPath
	fileInfo.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	t.Log.Infof("UploadFile fileInfo: %+v", fileInfo)
	if err := t.DB.Save(fileInfo).Error; err != nil {
		t.Log.Errorf("UploadFile DB err: %v", err)
		return "", err
	}
	if len(clis) > 0 {
		// 命令行调用
		fileUrl := t.Git.GetFileUrl(fileInfo.Path)
		t.Log.Infof("UploadFile success cli fileUrl: %s", fileUrl)
		go t.SyncDatabase()
		return fileUrl, nil
	}
	t.Log.Info("UploadFile success app")
	return "", nil
}

// ----------------------------------------------------------------

// DeleteFile 删除文件
func (t *App) DeleteFile(filePath string) *internal.Response {
	t.Log.Infof("DeleteFile filePath: %v", filePath)

	// 删除文件
	err := t.Git.Delete(filePath)
	if err != nil {
		t.Log.Errorf("DeleteFile Git err: %v", err)
	}

	// 更新 DB
	res := t.DB.Where("file_path", filePath).Delete(&internal.File{})
	if res.Error != nil {
		t.Log.Errorf("DeleteFile DB err: %v", res.Error)
		return internal.Fail(res.Error.Error())
	}

	go t.SyncDatabase()
	t.Log.Info("DeleteFile success")
	return internal.Success("操作成功")
}

// ----------------------------------------------------------------

// UpdateFileName 更新文件名称
func (t *App) UpdateFileName(filePath string, fileName string) *internal.Response {
	t.Log.Infof("UpdateFileName filePath: %v; fileName: %v", filePath, fileName)
	if fileName == "" {
		return internal.Fail("文件名称不能为空")
	}
	res := t.DB.Model(&internal.File{}).Where("file_path", filePath).Update("file_name", fileName)
	if res.Error != nil {
		t.Log.Errorf("UpdateFileName DB err: %v", res.Error)
		return internal.Fail(res.Error.Error())
	}
	t.Log.Info("UpdateFileName success")
	return internal.Success("操作成功")
}

// ----------------------------------------------------------------

// CopyFileUrl 复制链接到粘贴板
func (t *App) CopyFileUrl(fileUrl string) *internal.Response {
	t.Log.Info("CopyFileUrl fileUrl: ", fileUrl)
	err := clipboard.WriteAll(fileUrl)
	if err != nil {
		return internal.Fail(err.Error())
	}
	t.Log.Info("CopyFileUrl success")
	return internal.Success("已复制到粘贴板")
}

// ----------------------------------------------------------------

// SyncDatabase 同步本地数据库文件到远程仓库
func (t *App) SyncDatabase() {
	// 上传 cantor.db 文件
	err := t.Git.Update(config.GitDBFile, xfile.Read(t.DBFile))
	if err != nil {
		t.Log.Errorf("SyncDatabase err: %v", err)
	} else {
		t.Log.Info("SyncDatabase success")
	}
	return
}
