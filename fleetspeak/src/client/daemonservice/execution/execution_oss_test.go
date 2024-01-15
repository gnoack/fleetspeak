//go:build oss

package execution

import (
	"runtime"
	"testing"
)

func testClient(t *testing.T) string {
	if runtime.GOOS == "windows" {
		return `..\testclient\testclient.exe`
	}

	return "../testclient/testclient"
}
