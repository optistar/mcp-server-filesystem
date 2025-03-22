package top

import (
	"github.com/optistar/mcp-server-filesystem/tester"
	"testing"
)

func TestSearchFiles(t *testing.T) {
	tester.TestSearchFiles(tester.Wrap(t), tester.BypassFactory(Tools))
}
