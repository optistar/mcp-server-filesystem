package top

import (
	"github.com/optistar/mcp-server-filesystem/tester"
	"testing"
)

func TestCreateDirectory(t *testing.T) {
	tester.TestCreateDirectory(tester.Wrap(t), tester.BypassFactory(Tools))
}

func TestDirectoryTree(t *testing.T) {
	tester.TestDirectoryTree(tester.Wrap(t), tester.BypassFactory(Tools))
}

func TestEditFile(t *testing.T) {
	tester.TestEditFile(tester.Wrap(t), tester.BypassFactory(Tools))
}

func TestGetFileInfo(t *testing.T) {
	tester.TestGetFileInfo(tester.Wrap(t), tester.BypassFactory(Tools))
}

func TestListAllowedDirectories(t *testing.T) {
	tester.TestListAllowedDirectories(tester.Wrap(t), tester.BypassFactory(Tools))
}

func TestListDirectory(t *testing.T) {
	tester.TestListDirectory(tester.Wrap(t), tester.BypassFactory(Tools))
}

func TestMoveFile(t *testing.T) {
	tester.TestMoveFile(tester.Wrap(t), tester.BypassFactory(Tools))
}

func TestReadFile(t *testing.T) {
	tester.TestReadFile(tester.Wrap(t), tester.BypassFactory(Tools))
}

func TestReadMultipleFiles(t *testing.T) {
	tester.TestReadMultipleFiles(tester.Wrap(t), tester.BypassFactory(Tools))
}

func TestSearchFiles(t *testing.T) {
	tester.TestSearchFiles(tester.Wrap(t), tester.BypassFactory(Tools))
}

func TestWriteFile(t *testing.T) {
	tester.TestWriteFile(tester.Wrap(t), tester.BypassFactory(Tools))
}
