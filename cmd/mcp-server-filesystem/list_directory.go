package main

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"os"
	"strings"
)

func defineListDirectoryTool() mcp.Tool {
	return mcp.NewTool("list_directory",
		mcp.WithDescription(
			"Get a detailed listing of all files and directories in a specified path. "+
				"Results clearly distinguish between files and directories with [FILE] and [DIR] "+
				"prefixes. This tool is essential for understanding directory structure and "+
				"finding specific files within a directory. Only works within allowed directories."),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to the directory")),
	)
}

func listDirectoryHandler(_ context.Context, req mcp.CallToolRequest, allowedDirs []string) (*mcp.CallToolResult, error) {
	path, ok := req.Params.Arguments["path"].(string)
	if !ok {
		return mcp.NewToolResultError("path must be a string"), nil
	}
	validPath, err := validatePath(path, allowedDirs)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	entries, err := os.ReadDir(validPath)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if len(entries) == 0 {
		return mcp.NewToolResultText("Empty directory"), nil
	}
	var lines []string
	for _, entry := range entries {
		prefix := "[FILE]"
		if entry.IsDir() {
			prefix = "[DIR]"
		}
		lines = append(lines, fmt.Sprintf("%s %s", prefix, entry.Name()))
	}
	return mcp.NewToolResultText(strings.Join(lines, "\n")), nil
}
