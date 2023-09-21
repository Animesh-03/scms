package logger

import (
	"time"

	"github.com/fatih/color"
)

func LogConnectEvent(str string, a ...interface{}) {
	str = time.Now().Format(time.DateTime) + "|\t" + str
	// color.Green(str, a)
	color.New(color.FgGreen).Printf(str, a...)
}

func LogDisconnectEvent(str string, a ...interface{}) {
	str = time.Now().Format(time.DateTime) + "|\t" + str
	// color.Blue(str, a)
	color.New(color.FgBlue).Printf(str, a...)

}

func LogError(str string, a ...interface{}) {
	str = time.Now().Format(time.DateTime) + "|\t" + str
	// color.Red(str, a)
	color.New(color.FgRed).Printf(str, a...)

}

func LogInfo(str string, a ...interface{}) {
	str = time.Now().Format(time.DateTime) + "|\t" + str
	// color.HiMagenta(str, a)
	color.New(color.FgHiMagenta).Printf(str, a...)

}

func LogWarn(str string, a ...interface{}) {
	str = time.Now().Format(time.DateTime) + "|\t" + str
	// color.HiYellow(str, a)
	color.New(color.FgHiYellow).Printf(str, a...)

}
