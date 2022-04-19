package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/evercyan/brick/xfile"
	"github.com/evercyan/brick/xutil"
	"github.com/evercyan/cantor/backend/internal"
	"github.com/evercyan/cantor/config"
	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

// App ...
type App struct {
	Ctx     context.Context
	Log     *logrus.Logger
	Git     *internal.Git
	DB      *gorm.DB
	CfgFile string
	LogFile string
	DBFile  string
}

// NewApp ...
func NewApp() *App {
	return &App{}
}

// ----------------------------------------------------------------

// OnStartup 应用启动
func (t *App) OnStartup(ctx context.Context) {
	t.Ctx = ctx
	cfgPath := internal.GetCfgPath()
	// 日志
	t.LogFile = fmt.Sprintf(config.LogFile, cfgPath)
	t.Log = internal.NewLogger(t.LogFile)
	t.Log.Info("OnStartup begin")
	// Git
	t.CfgFile = fmt.Sprintf(config.CfgFile, cfgPath)
	t.Git = &internal.Git{}
	if err := yaml.Unmarshal([]byte(xfile.Read(t.CfgFile)), t.Git); err != nil {
		t.Log.Errorf("OnStartup cfgfile err: %v", err)
	}
	t.Log.Infof("OnStartup cfg: %+v", t.Git)
	// 数据库文件
	t.DBFile = fmt.Sprintf(config.DBFile, cfgPath)
	// 如果无 Git 配置, 不处理数据库相关
	if t.Git.Repo == "" {
		return
	}
	// 数据库处理
	t.database()
	return
}

// Menu 应用菜单
func (t *App) Menu() *menu.Menu {
	return menu.NewMenuFromItems(
		menu.SubMenu("Cantor", menu.NewMenuFromItems(
			menu.Text("关于 Cantor", nil, func(_ *menu.CallbackData) {
				t.diag(config.Description)
			}),
			menu.Text("检查更新", nil, func(_ *menu.CallbackData) {
				lastVersion := t.Git.GetLastVersion()
				needUpdate := config.Version < lastVersion
				msg := config.VersionNewMsg
				btns := []string{config.BtnConfirmText}
				if needUpdate {
					msg = fmt.Sprintf(config.VersionOldMsg, lastVersion)
					btns = []string{config.BtnConfirmText, config.BtnCancelText}
				}
				selection, err := t.diag(msg, btns...)
				if err != nil {
					return
				}
				if needUpdate && selection == config.BtnConfirmText {
					url := fmt.Sprintf(config.GitAppURL, lastVersion)
					runtime.BrowserOpenURL(t.Ctx, url)
				}
			}),
			menu.Separator(),
			menu.Text(
				"上传图片",
				keys.CmdOrCtrl("O"),
				func(_ *menu.CallbackData) {
					runtime.EventsEmit(t.Ctx, config.EventUploadBegin)
					resp := t.BatchUploadFile()
					if resp.Code == 0 {
						runtime.EventsEmit(t.Ctx, config.EventUploadSuccess, resp.Data)
					} else {
						runtime.EventsEmit(t.Ctx, config.EventUploadFail, resp.Msg)
					}
				},
			),
			menu.Separator(),
			menu.Text("退出", keys.CmdOrCtrl("Q"), func(_ *menu.CallbackData) {
				runtime.Quit(t.Ctx)
			}),
		)),
		menu.EditMenu(),
		menu.SubMenu("Help", menu.NewMenuFromItems(
			menu.Text(
				"打开配置文件",
				keys.Combo("C", keys.CmdOrCtrlKey, keys.ShiftKey),
				func(_ *menu.CallbackData) {
					if !xfile.IsExist(t.CfgFile) {
						t.diag("文件不存在, 请先更新配置")
						return
					}
					_, err := exec.Command("open", t.CfgFile).Output()
					if err != nil {
						t.diag("操作失败: " + err.Error())
						return
					}
				},
			),
			menu.Text(
				"打开日志文件",
				keys.Combo("L", keys.CmdOrCtrlKey, keys.ShiftKey),
				func(_ *menu.CallbackData) {
					if !xfile.IsExist(t.LogFile) {
						t.diag("文件不存在, 请先更新配置")
						return
					}
					_, err := exec.Command("open", t.LogFile).Output()
					if err != nil {
						t.diag("操作失败: " + err.Error())
						return
					}
				},
			),
			menu.Separator(),
			menu.Text(
				"打开应用主页",
				keys.Combo("H", keys.CmdOrCtrlKey, keys.ShiftKey),
				func(_ *menu.CallbackData) {
					runtime.BrowserOpenURL(t.Ctx, config.GitRepoURL)
				},
			),
		)),
	)
}

// OnDomReady ...
func (t *App) OnDomReady(ctx context.Context) {
	t.Log.Info("OnDomReady")
	return
}

// OnShutdown ...
func (t *App) OnShutdown(ctx context.Context) {
	t.Log.Info("OnShutdown")
	return
}

// OnBeforeClose ...
func (t *App) OnBeforeClose(ctx context.Context) bool {
	t.Log.Info("OnBeforeClose")
	// 返回 true 将阻止程序关闭
	return false
}

// ----------------------------------------------------------------

// migrate 数据同步
func (t *App) database() {
	t.Log.Info("OnStartup migrate begin")

	// 1. 如果 cantor.db 存在, 初始化, 返回
	if xfile.IsExist(t.DBFile) {
		t.DB = internal.NewDB(t.DBFile)
		return
	}

	// 2. 校验远程 cantor.db 是否存在, 存在直接同步后初始化, 返回
	dbContent := t.Git.GetContent(config.GitDBFile)
	if dbContent != "" {
		if err := xfile.Write(t.DBFile, dbContent); err != nil {
			t.Log.Errorf("OnStartup migrate sqlite err: %v", err)
			return
		}
		t.DB = internal.NewDB(t.DBFile)
		t.Log.Info("OnStartup migrate sqlite success")
		return
	}

	// 3. 校验远程 database.json 是否存在, 存在则迁移数据
	t.DB = internal.NewDB(t.DBFile)
	jsonContent := t.Git.GetContent("resource/database.json")
	if jsonContent != "" {
		list := make([]internal.File, 0)
		if err := json.Unmarshal([]byte(jsonContent), &list); err != nil {
			t.Log.Errorf("OnStartup migrate json err: %v", err)
			return
		}
		success := 0
		for i := len(list) - 1; i >= 0; i-- {
			res := t.DB.Create(&list[i])
			if res.Error == nil {
				success++
			}
		}
		t.Log.Infof("OnStartup migrate json success: %d", success)
	}

	return
}

// diag ...
func (t *App) diag(message string, buttons ...string) (string, error) {
	if len(buttons) == 0 {
		buttons = []string{
			config.BtnConfirmText,
		}
	}
	return runtime.MessageDialog(t.Ctx, runtime.MessageDialogOptions{
		Type:          runtime.InfoDialog,
		Title:         config.Title,
		Message:       message,
		CancelButton:  config.BtnConfirmText,
		DefaultButton: config.BtnConfirmText,
		Buttons:       buttons,
	})
}

// checkFile 校验上传文件
func (t *App) CheckFile(filepath string) error {
	if filepath == "" {
		return fmt.Errorf("请选择图片文件")
	}
	if t.Git.Repo == "" {
		return fmt.Errorf("请设置 Git 配置")
	}
	// 文件格式校验
	fileExt := strings.ToLower(path.Ext(filepath))
	if !xutil.IsContains(config.AllowFileExts, fileExt) {
		return fmt.Errorf("仅支持以下格式: %s", strings.Join(config.AllowFileExts, ", "))
	}
	// 文件大小校验
	fileSize := xfile.Size(filepath)
	t.Log.Infof("UploadFile fileSize: %v", fileSize)
	if fileSize > config.MaxFileSize {
		return fmt.Errorf("最大支持 4M 的文件")
	}
	return nil
}
