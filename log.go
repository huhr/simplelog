//wrapper.go
package simplelog

import ()

var (
	Global DomobLogger
)

func init() {
	Global = NewDomobLogger()
}

func LoadConfiguration(config Config) {
	Global.LoadConfiguration(config)
}

func Debug(format string, v ...interface{}) {
	Global.Debug(format, v ...)
}

func Info(format string, v ...interface{}) {
	Global.Info(format, v ...)
}

func Warn(format string, v ...interface{}) {
	Global.Warn(format, v ...)
}

func Error(format string, v ...interface{}) {
	Global.Error(format, v ...)
}

func Fatal(format string, v ...interface{}) {
	Global.Fatal(format, v ...)
}
