package main

import (
	"embed"
	"log"

	"github.com/evercyan/cantor/backend"
	"github.com/evercyan/cantor/config"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	app := backend.NewApp()
	err := wails.Run(&options.App{
		Title:             config.App,                           // 应用名称
		Width:             config.Width,                         // 初始宽度
		Height:            config.Height,                        // 初始高度
		MinWidth:          config.Width,                         // 最小宽度
		MinHeight:         config.Height,                        // 最小高度
		MaxWidth:          config.Width * 10,                    // 最大宽度
		MaxHeight:         config.Height * 10,                   // 最大高度
		DisableResize:     false,                                // 禁用调整窗口尺寸
		Frameless:         false,                                // 无边框
		StartHidden:       false,                                // 启动后即隐藏
		HideWindowOnClose: false,                                // 关闭窗口时将隐藏而不退出应用程序
		AssetServer:       &assetserver.Options{Assets: assets}, // 前端资源
		LogLevel:          logger.DEBUG,                         // 日志级别
		OnStartup:         app.OnStartup,                        // 程序启动回调
		OnDomReady:        app.OnDomReady,                       // 前端 dom 加载完成回调
		OnBeforeClose:     app.OnBeforeClose,                    // 关闭应用程序之前回调
		OnShutdown:        app.OnShutdown,                       // 程序退出回调
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
		Menu: app.Menu(),
	})
	if err != nil {
		log.Fatal(err)
	}
}
