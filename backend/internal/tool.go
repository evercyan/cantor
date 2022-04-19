package internal

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/evercyan/brick/xfile"
	"github.com/evercyan/cantor/config"
)

// GetCfgPath ...
func GetCfgPath() string {
	userPath, err := user.Current()
	if err != nil {
		panic("获取应用配置目录失败: " + err.Error())
	}
	cfgPath := fmt.Sprintf("%s/.%s", userPath.HomeDir, strings.ToLower(config.App))
	if !xfile.IsExist(cfgPath) {
		os.Mkdir(cfgPath, os.ModePerm)
	}
	return cfgPath
}
