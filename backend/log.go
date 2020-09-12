package backend

import (
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	logFile = "/tmp/cantor.log"
)

var log *logrus.Logger
var once sync.Once

// Log ...
func Log() *logrus.Logger {
	once.Do(func() {
		log = logrus.New()
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			panic("创建日志文件失败: " + err.Error())
		}
		mw := io.MultiWriter(os.Stdout, file)
		log.SetOutput(mw)
		log.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
	})
	return log
}
