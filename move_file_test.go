package top

import (
	"github.com/optistar/mcp-server-filesystem/tester"
	"testing"
)

func TestMoveFile(t *testing.T) {
	tester.TestMoveFile(tester.Wrap(t), tester.BypassFactory(Tools))
}
