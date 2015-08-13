//log_test.go
package simplelog

import (
	"testing"
)

func TestLog(t *testing.T) {
	LoadConfiguration("log.json")
	Debug("Hello Info")
	Error("Hello Error")
}
