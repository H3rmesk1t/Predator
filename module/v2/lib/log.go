package lib

import (
	"Predator/pkg/config"
	"github.com/kataras/golog"
	"io"
	"os"
)

var Log = golog.New()

func InitLog() {
	if config.SpyDebug {
		Log.SetLevel("debug")
	}
	if config.SpySilent {
		Log.SetLevel("error")
		Log.SetTimeFormat("")
	}
	logFile, err := os.OpenFile("spy.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	Log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}
