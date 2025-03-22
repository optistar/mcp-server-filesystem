package main

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"os"
	"strings"
)

func defineReadMultipleFilesTool() mcp.Tool {
	return mcp.NewTool("read_multiple_files",
		mcp.WithDescription(
			"Read the contents of multiple files simultaneously. This is more "+
				"efficient than reading files one by one when you need to analyze "+
				"or compare multiple files. Each file's content is returned with its "+
				"path as a reference. Failed reads for individual files won't stop "+
				"the entire operation. Only works within allowed directories."),
		mcp.WithArray("paths",
			mcp.Required(),
			mcp.Description("Array of file paths"),
			mcp.Items(map[string]interface{}{
				"type": "string",
			})),
	)
}

func readMultipleFilesHandler(_ context.Context, req mcp.CallToolRequest, allowedDirs []string) (*mcp.CallToolResult, error) {
	paths, ok := req.Params.Arguments["paths"].([]interface{})
	if !ok {
		return mcp.NewToolResultError("paths must be an array"), nil
	}
	var results []string
	for _, p := range paths {
		path, ok := p.(string)
		if !ok {
			results = append(results, fmt.Sprintf("%v: Error - must be a string", p))
			continue
		}
		validPath, err := validatePath(path, allowedDirs)
		if err != nil {
			results = append(results, fmt.Sprintf("%s: Error - %v", path, err))
			continue
		}
		content, err := os.ReadFile(validPath)
		if err != nil {
			results = append(results, fmt.Sprintf("%s: Error - %v", path, err))
			continue
		}
		results = append(results, fmt.Sprintf("%s:\n%s", path, string(content)))
	}
	return mcp.NewToolResultText(strings.Join(results, "\n---\n")), nil
}
