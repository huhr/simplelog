package simplelog

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)


var VERSION = "0.2.0"

var (
	global = make(map[string]*SimpleLogger, 5)
)

func Debug(format string, v ...interface{}) {
	global["root"].Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	global["root"].Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	global["root"].Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	global["root"].Error(format, v...)
}

func Fatal(format string, v ...interface{}) {
	global["root"].Fatal(format, v...)
}

// load configration from json file
func LoadConfigFile(filename string) error {
	cfgs, err := readConfiguration(filename)
	if err != nil {
		return err
	}
	return LoadConfigMap(cfgs)
}

// support config map
func LoadConfigMap(cfgs map[string][]map[string]string) error {
	if _, ok := cfgs["root"]; !ok {
		return errors.New("log is missing")
	}
	for k, v := range cfgs {
		logger, err := NewSimpleLogger(k, v)
		if err != nil {
			return err
		}
		global[k] = logger
	}
	return nil
}

// read config file into map
func readConfiguration(filename string) (cfgs map[string][]map[string]string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	src, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(src, &cfgs)
	if err != nil {
		return nil, err
	}
	return cfgs, nil
}

// 根据名称获取一个logger实例
func GetLogger(name string) (*SimpleLogger, error) {
	if name == "root" {
		return nil, errors.New("can not get root logger, just use static log")
	}
	if logger, ok := global[name]; ok {
		return logger, nil
	}
	return nil, errors.New(fmt.Sprintf("logger %s is not exits", name))
}
