//auther  huhaoran<huhaoran@domob.com>
//log.go
package simplelog

import (
	"time"
	"fmt"
	"runtime"
	"path"
	"os"
	"sync"
	"io"
)


//define format
const (
	//[log_level] [time] [filename]:[line] [message]
	TEXT_FORMAT = "%s %s %s:%d :: %s \n"
	//[message]
	BRIFE_FORMAT = "%s \n"
)

//define level
const (
	LEVEL_DEBUG=iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
)

//define level string
var LEVEL_STRINGS = [...]string{
	"debug",
	"info",
	"warn",
	"error",
	"fatal",
}

//DomobLogger just for golang
type DomobLogger struct {
	mu         sync.Mutex
	//every level can hold it`s out
	out        map[int]io.Writer
	level      int
	//every level can hold it`s format
	//default os.Stdout
	format     map[int]string
}

func (l *DomobLogger)SetLevel(level int) {
	l.level = level
}

func NewDomobLogger() DomobLogger{
	l := DomobLogger{}
	//Debug, Info, Warn, Error, Fatal
	l.format = make(map[int]string, 5)
	l.out = make(map[int]io.Writer, 5)
	return l
}

func (l *DomobLogger)Output(level int, message string) {
	if level < l.level {
		return
	}

	// first select the Format of the log
	var log string
	switch l.format[level] {
	case TEXT_FORMAT:
		time_str := fmt.Sprintf("%s", time.Now())[:19]
		// skip three
		// here we get the filename so we need not to name the logger
		_, filename, line, _ := runtime.Caller(3)
		_, filename = path.Split(filename)
		log = fmt.Sprintf(TEXT_FORMAT, LEVEL_STRINGS[level], time_str, filename, line, message)
	case BRIFE_FORMAT:
		log = fmt.Sprintf(BRIFE_FORMAT, message)
	default:
		panic("there is no format")
	}

	// lock the func
	l.mu.Lock()
	defer l.mu.Unlock()
	_ ,err := l.out[level].Write([]byte(log))
	if err !=nil {
		fmt.Println(err.Error())
	}
}

func (l *DomobLogger) Debug(format string, v ...interface{}) {
	l.Output(LEVEL_DEBUG, fmt.Sprintf(format, v ...))
}

func (l *DomobLogger) Info(format string, v ...interface{}) {
	l.Output(LEVEL_INFO, fmt.Sprintf(format, v ...))
}

func (l *DomobLogger) Warn(format string, v ...interface{}) {
	l.Output(LEVEL_WARN, fmt.Sprintf(format, v ...))
}

func (l *DomobLogger) Error(format string, v ...interface{}) {
	l.Output(LEVEL_ERROR, fmt.Sprintf(format, v ...))
}

func (l *DomobLogger) Fatal(format string, v ...interface{}) {
	l.Output(LEVEL_FATAL, fmt.Sprintf(format, v ...))
}

// load configration from .cfg file
func (l *DomobLogger) LoadConfiguration(config Config) {
	//first set level
	switch config.Basic.Level {
	case "DEBUG":
		l.SetLevel(LEVEL_DEBUG)
	case "INFO":
		l.SetLevel(LEVEL_INFO)
	case "WARN":
		l.SetLevel(LEVEL_WARN)
	case "ERROR":
		l.SetLevel(LEVEL_ERROR)
	case "FATAL":
		l.SetLevel(LEVEL_FATAL)
	default:
		//this may print out the screen use the TEXT_FORMAT or >> to he console file
		l.Info("default log level DEBUG")
		l.SetLevel(LEVEL_DEBUG)
	}

	//second set format for each level
	for k, format := range []string{config.Debug.Format, config.Info.Format, config.Warn.Format, config.Error.Format, config.Fatal.Format} {
		switch format {
		case "TEXT":
			l.format[k] = TEXT_FORMAT
		case "BRIFE":
			l.format[k] = BRIFE_FORMAT
		default:
			l.format[k] = TEXT_FORMAT
		}
	}

	//then set io.Writer for each level
	for k, file := range []string{config.Debug.File, config.Info.File, config.Warn.File, config.Error.File, config.Fatal.File} {
		if file != "" {
			if _, err := os.Stat(file); os.IsNotExist(err) {
				l.out[k], err = os.Create(file)
				if err != nil {
					panic(err.Error())
				}
			}else {
				l.out[k], err = os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0600)
				if err != nil {
					panic(err.Error())
				}
			}
		}else {
			l.out[k] = os.Stdout
		}
	}

	//finally start go goroutinue to cut the log
	for k, hourly := range [][]bool{[]bool{config.Debug.Hourly, config.Debug.Daily}, []bool{config.Info.Hourly, config.Info.Daily},
		[]bool{config.Warn.Hourly, config.Warn.Daily}, []bool{config.Error.Hourly, config.Error.Daily},
		[]bool{config.Fatal.Hourly, config.Fatal.Daily}} {

		if hourly[0] {
			go l.cut_log(config, k, true)
		} else if hourly[1] {
			go l.cut_log(config, k, false)
		}else {
			l.Output(k, "log file will not be cutted")
		}
	}
}

// Basic中保存缺省的配置和日志等级的控制
// 每个级别单独配置输出文件，format和切割控制
type Config struct {
	Basic struct {
		Level     string
	}

	Info struct {
		File      string
		Hourly    bool
		Daily     bool
		Format   string
	}

	Debug struct {
		File      string
		Hourly    bool
		Daily     bool
		Format   string
	}

	Warn struct {
		File      string
		Hourly    bool
		Daily     bool
		Format   string
	}

	Error struct {
		File      string
		Hourly    bool
		Daily     bool
		Format   string
	}

	Fatal struct {
		File      string
		Hourly    bool
		Daily     bool
		Format   string
	}
}

func NewConfig(file string, hourly, daily bool, level, format string)  Config{
	config := Config{}
	config.Basic.Level = level
	config.Debug.File = file
	config.Debug.Hourly = hourly
	config.Debug.Daily = daily
	config.Debug.Format = format
	return config
}

//cut_log hourly or daily
func (l *DomobLogger)cut_log(config Config, level int, hourly bool) {
	//sleep
	time.Sleep(time.Duration(getInterval(hourly))*time.Second)

	//wake up to cut the log file
	files := []string{config.Debug.File, config.Info.File, config.Warn.File, config.Error.File, config.Fatal.File}
	name := fmt.Sprintf("%s.%s", files[level], time.Now().Format("20060102"))
	os.Rename(files[level], name)
	l.out[level], _ = os.Create(files[level])
}

//get the sleep time interval
func getInterval(hourly bool) int64 {
	now := time.Now()
	if hourly {
		cut_time := time.Unix(time.Now().Unix() + 3600, 0).Round(time.Hour)
		return cut_time.Unix() - now.Unix() + 1
	}else {
		y,m,d := now.AddDate(0,0,1).Date()
		cut_time := time.Date(y,m,d, 0, 0, 1, 0, time.Local)
		return cut_time.Unix() - now.Unix() + 1
	}
}
