package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true, // 显示完整时间戳
		TimestampFormat: "2006-01-02 15:04:05",
	})
	Logger.SetLevel(logrus.DebugLevel)
	Logger.SetOutput(os.Stdout)
}
