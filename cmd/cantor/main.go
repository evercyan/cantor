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

var (
	flagCli = true
)

func main() {
	root := &cobra.Command{
		Use:     "cantor",
		Short:   "cantor: 一个简单好用的图床应用",
		Version: "v0.0.1",
	}
	root.AddCommand(&cobra.Command{
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

			for _, filepath := range args {
				if !xfile.IsExist(filepath) || !xfile.IsFile(filepath) {
					xcolor.Fail(filepath, "无效的图片路径")
					continue
				}
				if err := app.CheckFile(filepath); err != nil {
					xcolor.Fail(filepath, err.Error())
					continue
				}
				fileUrl, err := app.Upload(filepath, true)
				if err != nil {
					xcolor.Fail(filepath, err.Error())
					continue
				}
				if flagCli {
					xcolor.Success(filepath, fileUrl)
				} else {
					fmt.Println(fileUrl)
				}
			}
		},
	})
	root.PersistentFlags().BoolVarP(
		&flagCli, "cli", "c", true, "run mode, default cli",
	)
	root.Execute()
}
