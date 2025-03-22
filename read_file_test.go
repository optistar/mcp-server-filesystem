package top

import (
	"github.com/optistar/mcp-server-filesystem/tester"
	"testing"
)

func TestReadFile(t *testing.T) {
	tester.TestReadFile(tester.Wrap(t), tester.BypassFactory(Tools))
}
