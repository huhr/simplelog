//log_test.go
package simplelog

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	config := NewConfig("/home/huhaoran/explore/simplelog/out.log", false, true, "DEBUG", "TEXT")
	LoadConfiguration(config)
	go func() {
		for i:=0;i<10;i++ {
			Debug("hello simplelog. logong asdasdqwaclkjlkasjdlm,zxncqoiweoqwpeoqiwe,smnc,kljdlkaoiqwueoiqwejkasnand,sa")
		}
	}()

	time.Sleep(1)
	go func() {
		for i:=0;i<10;i++ {
			Info("hello simplelog. logong asdasdqwaclkjlkasjdlm,zxncqoiweoqwpeoqiwe,smnc,kljdlkaoiqwueoiqwejkasnand,sa")
		}
	}()
}
