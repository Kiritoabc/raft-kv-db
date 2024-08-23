package log

import (
	log "github.com/sirupsen/logrus"
	"os"
)

var Log = log.New()

func init() {
	// 为当前logrus实例设置消息的输出，同样地，
	// 可以设置logrus实例的输出到任意io.writer
	Log.Out = os.Stdout

	// 为当前logrus实例设置消息输出格式为json格式。
	// 同样地，也可以单独为某个logrus实例设置日志级别和hook，这里不详细叙述。
	Log.Formatter = &log.JSONFormatter{}

	// 设置日志级别
	Log.SetLevel(log.InfoLevel)
	Log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	// 调用方法添加字段
	Log.SetReportCaller(true)

	Log.Info("Log init success...")
}
