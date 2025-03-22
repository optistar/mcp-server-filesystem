package top

import (
	"github.com/optistar/mcp-server-filesystem/tester"
	"testing"
)

func TestCreateDirectory(t *testing.T) {
	tester.TestCreateDirectory(tester.Wrap(t), tester.BypassFactory(Tools))
}
