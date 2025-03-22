package top

import (
	"github.com/optistar/mcp-server-filesystem/tester"
	"testing"
)

func TestGetFileInfo(t *testing.T) {
	tester.TestGetFileInfo(tester.Wrap(t), tester.BypassFactory(Tools))
}
