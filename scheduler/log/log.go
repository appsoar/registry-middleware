package log

import (
	//	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
)

var (
	Logger     *logs.BeeLogger
	debug      = false
	LogChannel = make(chan string)
)

type UserLog struct {
	user string
}

type WebLog struct {
	Time   string `json:"time"`
	User   string `json:"user"`
	Level  string `json:"level"`
	Detail string `json:"detail"`
}

func Debug(format string, v ...interface{}) {
	Logger.Debug(format, v)
	go func() {
		var str string
		if len(v) != 0 {
			str = fmt.Sprintf(format, v)
		} else {
			str = format
		}

		LogChannel <- str
	}()
}

func init() {
	if os.Getenv("DEBUG") == "true" {
		debug = true

	}
	Logger = logs.NewLogger(10000)
	if debug {
		Logger.SetLevel(logs.LevelDebug) //设置写到缓冲区的日志
		Logger.EnableFuncCallDepth(true) //输出时显示文件名和行号
	} else {
		Logger.SetLevel(logs.LevelInfo)
	}
	Logger.SetLogger("console", "")
	//	以下需要设置接口
	//	Logger.SetLogger("file", `{"filename":"test.log"}`)
	//	Logger.SetLogger("smtp", `{"username":"kiongf@126.com","password":    "wangweihong1988","fromAddress":"kiongf@126.com","host":"smtp.126.c    om:25","sendTos":["wwhvw@126.com"] }`)
}
