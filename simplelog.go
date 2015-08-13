// Auther huhaoran<huhaoran@domob.com>
package simplelog

import (
	"strings"
	"fmt"
	"io"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

// Define format
const (
	// [log_level] [time] [filename]:[line] [message]
	DETAIL_FORMAT = "%s %s %s:%d :: %s \n"
	// [message]
	BRIFE_FORMAT = "%s \n"
)

// Define level
const (
	LEVEL_DEBUG = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
)

// Define level string
var LEVEL_STRINGS = [...]string{
	"debug",
	"info",
	"warn",
	"error",
	"fatal",
}

// SimpleLogger just for golang
type SimpleLogger struct {
	// every level can hold it`s out
	writers   map[int]*LogWriter
}

func NewSimpleLogger() SimpleLogger {
	l := SimpleLogger{}
	// Debug, Info, Warn, Error, Fatal
	l.writers = make(map[int]*LogWriter, 0)
	return l
}

func (l *SimpleLogger) Debug(format string, v ...interface{}) {
	if writer, ok := l.writers[LEVEL_DEBUG]; ok {
		writer.Output(LEVEL_DEBUG, fmt.Sprintf(format, v...))
	}
}

func (l *SimpleLogger) Info(format string, v ...interface{}) {
	if writer, ok := l.writers[LEVEL_DEBUG]; ok {
		writer.Output(LEVEL_INFO, fmt.Sprintf(format, v...))
	}
}

func (l *SimpleLogger) Warn(format string, v ...interface{}) {
	if writer, ok := l.writers[LEVEL_WARN]; ok {
		writer.Output(LEVEL_WARN, fmt.Sprintf(format, v...))
	}
}

func (l *SimpleLogger) Error(format string, v ...interface{}) {
	if writer, ok := l.writers[LEVEL_ERROR]; ok {
		writer.Output(LEVEL_ERROR, fmt.Sprintf(format, v...))
	}
}

func (l *SimpleLogger) Fatal(format string, v ...interface{}) {
	if writer, ok := l.writers[LEVEL_FATAL]; ok {
		writer.Output(LEVEL_FATAL, fmt.Sprintf(format, v...))
	}
}

// load configration from json file
func (l *SimpleLogger) LoadConfiguration(filename string) {
	cfgs := l.readConfiguration(filename)
	for _, cfg := range cfgs {
		writer := NewLogWriter(cfg["Out"], cfg["Format"], cfg["Cut"])
		for levelIndex, levelStr := range LEVEL_STRINGS {
			if strings.Contains(cfg["Level"], levelStr) {
				l.writers[levelIndex] = writer
			}
		}
	}
}

// read config file into map 
func (l SimpleLogger) readConfiguration(filename string) (cfg []map[string]string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	src, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err.Error())
	}
	err = json.Unmarshal(src, &cfg)
	if err != nil {
		panic(err.Error())
	}
	return cfg
}

// LogWriter
type LogWriter struct {
	sync.Mutex
	filename  string
	format    string
	cutType   string
	writer    io.Writer
}

func NewLogWriter(filename, format, cutType string) *LogWriter {
	var writer io.Writer
	switch filename {
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	default:
		writer = func() (w io.Writer) {
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				w, err = os.Create(filename)
				if err != nil {
					panic(err.Error())
				}
			} else {
				w, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
				if err != nil {
					panic(err.Error())
				}
			}
			return
		}()
	}
	w := &LogWriter{sync.Mutex{},filename, format, cutType, writer}
	go w.cutLog()
	return w
}

func (w *LogWriter) Output(level int, msg string) {
	var log string
	switch w.format {
	case "detail":
		time_str := fmt.Sprintf("%s", time.Now())[:19]
		// skip three
		// here we get the filename so we need not to name the logger
		_, filename, line, _ := runtime.Caller(3)
		_, filename = path.Split(filename)
		log = fmt.Sprintf(DETAIL_FORMAT, LEVEL_STRINGS[level], time_str, filename, line, msg)
	case "brife":
		log = fmt.Sprintf(BRIFE_FORMAT, msg)
	default:
		panic("there is no format")
	}
	// lock the func
	w.Lock()
	defer w.Unlock()
	_, err := w.writer.Write([]byte(log))
	if err != nil {
		fmt.Println(err.Error())
	}
}

// cut by hourly or daily
func (w *LogWriter) cutLog() {
	for {
		var (
			t int64
			name string
		)
		now := time.Now()
		switch w.cutType {
		case "daily":
			t = time.Unix(time.Now().Unix() + 3600, 0).Round(time.Hour).Unix() - now.Unix() + 1
			name = fmt.Sprintf("%s.%s", w.filename, time.Now().Format("20060102"))
		case "hourly":
			y, m, d := now.AddDate(0, 0, 1).Date()
			t = time.Date(y, m, d, 0, 0, 1, 0, time.Local).Unix() - now.Unix() + 1
			name = fmt.Sprintf("%s.%s", w.filename, time.Now().Format("2006010215"))
		default:
			return
		}
		// sleep
		time.Sleep(time.Duration(t) * time.Second)

		// wake up and mv the log file
		os.Rename(w.filename, name)
		writer, err := os.Create(w.filename)
		if err != nil {
			fmt.Printf("create new log file error: %s", err.Error())
			continue
		}
		w.writer = writer
	}
}
