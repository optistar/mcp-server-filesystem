package top

import (
	"github.com/optistar/mcp-server-filesystem/tester"
	"testing"
)

func TestListAllowedDirectories(t *testing.T) {
	tester.TestListAllowedDirectories(tester.Wrap(t), tester.BypassFactory(Tools))
}
