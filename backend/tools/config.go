package tools

import (
	"fmt"
	"os"
	"os/user"

	"github.com/evercyan/cantor/backend/configs"
	"github.com/evercyan/letitgo/util"
)

// GetConfigPath ...
func GetConfigPath() string {
	user, err := user.Current()
	if err != nil {
		panic("获取应用配置目录失败: " + err.Error())
	}
	configPath := fmt.Sprintf("%s/.%s", user.HomeDir, configs.AppName)
	if !util.IsExist(configPath) {
		os.Mkdir(configPath, os.ModePerm)
	}
	return configPath
}

// GetLogFilePath ...
func GetLogFilePath() string {
	return fmt.Sprintf("%s/%s.log", GetConfigPath(), configs.AppName)
}
