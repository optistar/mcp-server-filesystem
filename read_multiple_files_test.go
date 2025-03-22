package top

import (
	"github.com/optistar/mcp-server-filesystem/tester"
	"testing"
)

func TestReadMultipleFiles(t *testing.T) {
	tester.TestReadMultipleFiles(tester.Wrap(t), tester.BypassFactory(Tools))
}
