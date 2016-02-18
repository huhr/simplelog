// @auther huhaoran<huhaoran@domob.com>
package simplelog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

// define format
const (
	// [time] [log_level] [log_name] [caller] :: [message]
	DETAIL_FORMAT = "%s %s %s %s:%d :: %s \n"
	// [message]
	BRIFE_FORMAT = "%s \n"
)

// define log level
const (
	LEVEL_DEBUG = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
)

// define level string
var LEVEL_STRINGS = [...]string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"FATAL",
}

// logWriters map
var logWriters = make(map[string]*LogWriter)

// SimpleLogger just for golang
type SimpleLogger struct {
	// logger name
	name string
	// every level has its format
	formats map[int]string
	writers map[int]*LogWriter
}

func NewSimpleLogger(name string, cfgs []map[string]string) (*SimpleLogger, error) {
	l := &SimpleLogger{name: name}

	l.writers = make(map[int]*LogWriter, 5)
	l.formats = make(map[int]string, 5)
	for _, cfg := range cfgs {
		var level, output, rotation, format string
		var exists bool
		if level, exists = cfg["Level"]; !exists {
			return nil, errors.New("level section is missing")
		}
		if output, exists = cfg["Output"]; !exists {
			return nil, errors.New("Output section is missing")
		}
		if rotation, exists = cfg["Rotation"]; !exists {
			return nil, errors.New("Rotation section is missing")
		}
		if format, exists = cfg["Format"]; !exists {
			return nil, errors.New("format section is missing")
		}

		// one output one LogWriter, since file modification is unsafe in concurrency
		var writer *LogWriter
		if w, ok := logWriters[output]; ok {
			writer = w
		} else {
			writer = NewLogWriter(output, rotation)
			logWriters[output] = writer
		}

		for levelIndex, levelStr := range LEVEL_STRINGS {
			if strings.Contains(level, strings.ToLower(levelStr)) {
				// this level  have been initialized at previous nodes
				if _, ok := l.formats[levelIndex]; ok {
					continue
				}
				l.formats[levelIndex] = format
				l.writers[levelIndex] = writer
			}
		}
	}
	return l, nil
}

func (l *SimpleLogger) Debug(f string, v ...interface{}) {
	if writer, ok := l.writers[LEVEL_DEBUG]; ok {
		writer.Output(LEVEL_DEBUG, l.name, l.formats[LEVEL_DEBUG], fmt.Sprintf(f, v...))
	}
}

func (l *SimpleLogger) Info(f string, v ...interface{}) {
	if writer, ok := l.writers[LEVEL_DEBUG]; ok {
		writer.Output(LEVEL_INFO, l.name, l.formats[LEVEL_INFO], fmt.Sprintf(f, v...))
	}
}

func (l *SimpleLogger) Warn(f string, v ...interface{}) {
	if writer, ok := l.writers[LEVEL_WARN]; ok {
		writer.Output(LEVEL_WARN, l.name, l.formats[LEVEL_WARN], fmt.Sprintf(f, v...))
	}
}

func (l *SimpleLogger) Error(f string, v ...interface{}) {
	if writer, ok := l.writers[LEVEL_ERROR]; ok {
		writer.Output(LEVEL_ERROR, l.name, l.formats[LEVEL_ERROR], fmt.Sprintf(f, v...))
	}
}

func (l *SimpleLogger) Fatal(f string, v ...interface{}) {
	if writer, ok := l.writers[LEVEL_FATAL]; ok {
		writer.Output(LEVEL_FATAL, l.name, l.formats[LEVEL_FATAL], fmt.Sprintf(f, v...))
	}
}

// LogWriter manage Writer, Format & Rotation
type LogWriter struct {
	sync.Mutex
	output   string
	rotation string
	writer   io.Writer
}

func NewLogWriter(output, rotation string) *LogWriter {
	var writer io.Writer
	switch output {
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	default:
		writer = func() (w io.Writer) {
			if _, err := os.Stat(output); os.IsNotExist(err) {
				w, err = os.Create(output)
				if err != nil {
					panic(err.Error())
				}
			} else {
				w, err = os.OpenFile(output, os.O_APPEND|os.O_WRONLY, 0600)
				if err != nil {
					panic(err.Error())
				}
			}
			return
		}()
	}
	w := &LogWriter{sync.Mutex{}, output, rotation, writer}
	go w.rotate()
	return w
}

// output msg to io.Writer
func (w *LogWriter) Output(level int, loggerName, format, msg string) {
	var log string
	switch format {
	case "detail":
		timeStr := fmt.Sprintf("%s", time.Now())[:19]
		skip := 2
		if loggerName == "root" {
			skip = 3
		}
		_, filename, line, _ := runtime.Caller(skip)
		_, filename = path.Split(filename)
		log = fmt.Sprintf(DETAIL_FORMAT, timeStr,  LEVEL_STRINGS[level], loggerName, filename, line, msg)
	case "brife":
		log = fmt.Sprintf(BRIFE_FORMAT, msg)
	default:
		log = fmt.Sprintf(BRIFE_FORMAT, msg)
	}
	// lock the func
	w.Lock()
	defer w.Unlock()
	_, err := w.writer.Write([]byte(log))
	if err != nil {
		fmt.Println(err.Error())
	}
}

// rotate by hourly or daily
func (w *LogWriter) rotate() {
	for {
		var (
			t    int64
			name string
		)
		now := time.Now()
		switch w.rotation {
		case "hourly":
			t = time.Unix(time.Now().Unix()+3600, 0).Round(time.Hour).Unix() - now.Unix() + 1
			name = fmt.Sprintf("%s.%s", w.output, time.Now().Format("2006010215"))
		case "daily":
			y, m, d := now.AddDate(0, 0, 1).Date()
			t = time.Date(y, m, d, 0, 0, 1, 0, time.Local).Unix() - now.Unix() + 1
			name = fmt.Sprintf("%s.%s", w.output, time.Now().Format("20060102"))
		default:
			return
		}
		// sleep
		time.Sleep(time.Duration(t) * time.Second)

		// wake up and mv the log file
		os.Rename(w.output, name)
		writer, err := os.Create(w.output)
		if err != nil {
			fmt.Printf("create new log file error: %s", err.Error())
			continue
		}
		w.writer = writer
	}
}

