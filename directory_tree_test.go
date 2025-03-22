package top

import (
	"github.com/optistar/mcp-server-filesystem/tester"
	"testing"
)

func TestDirectoryTree(t *testing.T) {
	tester.TestDirectoryTree(tester.Wrap(t), tester.BypassFactory(Tools))
}
