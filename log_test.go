//log_test.go
package simplelog

import (
	"runtime"
	"fmt"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	err := LoadConfigFile("config.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	logger, err := GetLogger("another")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	d := func() {
		i := 1
		for i <= 100 {
			logger.Debug("Hello Info")
			i++
		}
	}
	go d()
	time.Sleep(2 * time.Second)
}
