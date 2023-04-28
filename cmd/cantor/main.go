package main

import (
	"context"
	"fmt"
	"os"

	"github.com/evercyan/brick/xcli/xcolor"
	"github.com/evercyan/brick/xfile"
	"github.com/evercyan/cantor/backend"
	"github.com/spf13/cobra"
)

// env 运行环境
var env string

// ...
const (
	envTypora = "typora"
)

// ...
var (
	root = &cobra.Command{
		Use:     "cantor",
		Short:   "一个简单好用的图床应用",
		Version: "v0.1.0",
	}
	upload = &cobra.Command{
		Use:   "upload",
		Short: "批量上传图片文件, e.g. cantor upload ~/demo.png ~/demo2.png",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				xcolor.Fail("Error:", "图片路径不能为空")
				return
			}
			app := backend.NewApp()
			app.OnStartup(context.Background())
			if app.Git.Repo == "" {
				xcolor.Fail("Error:", "无效的配置文件, 请先安装 Cantor 应用设置 Git 配置")
				os.Exit(0)
			}
			for _, fpath := range args {
				if !xfile.IsExist(fpath) || !xfile.IsImage(fpath) {
					xcolor.Fail("✘", "无效的图片")
					continue
				}
				if err := app.CheckFile(fpath); err != nil {
					xcolor.Fail("✘", err.Error())
					continue
				}
				fileUrl, err := app.Upload(fpath, true)
				if err != nil {
					xcolor.Fail("✘", err.Error())
					continue
				}
				if env == envTypora {
					// 在 typora 图片上传使用时只输出文件路径
					fmt.Println(fileUrl)
				} else {
					xcolor.Success("➤", fpath)
					xcolor.Success("✔︎", fileUrl)
				}
			}
		},
	}
)

func init() {
	root.AddCommand(upload)
	root.PersistentFlags().StringVarP(&env, "env", "", "cli", "env, default cli")
}

func main() {
	root.Execute()
}
