package top

import (
	"github.com/optistar/mcp-server-filesystem/tester"
)

// Define tools
var Tools = []tester.ToolHandler{
	{DefineReadFileTool(), ReadFileHandler},
	{DefineReadMultipleFilesTool(), ReadMultipleFilesHandler},
	{DefineWriteFileTool(), WriteFileHandler},
	{DefineEditFileTool(), EditFileHandler},
	{DefineCreateDirectoryTool(), CreateDirectoryHandler},
	{DefineListDirectoryTool(), ListDirectoryHandler},
	{DefineDirectoryTreeTool(), DirectoryTreeHandler},
	{DefineMoveFileTool(), MoveFileHandler},
	{DefineSearchFilesTool(), SearchFilesHandler},
	{DefineGetFileInfoTool(), GetFileInfoHandler},
	{DefineListAllowedDirectoriesTool(), ListAllowedDirectoriesHandler},
}
