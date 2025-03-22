package top

import (
	"github.com/optistar/mcp-server-filesystem/tester"
	"testing"
)

func TestListDirectory(t *testing.T) {
	tester.TestListDirectory(tester.Wrap(t), tester.BypassFactory(Tools))
}
