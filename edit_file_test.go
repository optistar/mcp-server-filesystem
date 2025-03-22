package top

import (
	"github.com/optistar/mcp-server-filesystem/tester"
	"testing"
)

func TestEditFile(t *testing.T) {
	tester.TestEditFile(tester.Wrap(t), tester.BypassFactory(Tools))
}
