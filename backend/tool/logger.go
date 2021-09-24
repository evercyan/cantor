package tool

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"sync"

	"github.com/evercyan/cantor/backend/cfg"
	"github.com/evercyan/letitgo/file"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
	once   sync.Once
)

// NewLogger ...
func NewLogger() *logrus.Logger {
	once.Do(func() {
		logger = logrus.New()
		filePath := fmt.Sprintf(cfg.LogFile, GetConfigPath())
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			panic("创建日志文件失败: " + err.Error())
		}
		logger.SetOutput(io.MultiWriter(os.Stdout, f))
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
	})
	return logger
}

// GetConfigPath ...
func GetConfigPath() string {
	userPath, err := user.Current()
	if err != nil {
		panic("获取应用配置目录失败: " + err.Error())
	}
	configPath := fmt.Sprintf("%s/.%s", userPath.HomeDir, cfg.AppName)
	if !file.IsExist(configPath) {
		os.Mkdir(configPath, os.ModePerm)
	}
	return configPath
}
