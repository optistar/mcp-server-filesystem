package top

import (
	"github.com/optistar/mcp-server-filesystem/tester"
	"testing"
)

func TestWriteFile(t *testing.T) {
	tester.TestWriteFile(tester.Wrap(t), tester.BypassFactory(Tools))
}
