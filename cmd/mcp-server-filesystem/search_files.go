package main

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	"os"
	"path/filepath"
	"strings"
)

func defineSearchFilesTool() mcp.Tool {
	return mcp.NewTool("search_files",
		mcp.WithDescription(
			"Recursively search for files and directories matching a pattern. "+
				"Searches through all subdirectories from the starting path. The search "+
				"is case-insensitive and matches partial names. Returns full paths to all "+
				"matching items. Great for finding files when you don't know their exact location. "+
				"Only searches within allowed directories."),
		mcp.WithString("path", mcp.Required(), mcp.Description("Starting path")),
		mcp.WithString("pattern", mcp.Required(), mcp.Description("Search pattern")),
		mcp.WithArray("excludePatterns",
			mcp.Description("Patterns to exclude"),
			func(schema map[string]interface{}) {
				schema["default"] = []interface{}{}
			},
			mcp.Items(map[string]interface{}{
				"type": "string",
			}),
		),
	)
}

func searchFilesHandler(_ context.Context, req mcp.CallToolRequest, allowedDirs []string) (*mcp.CallToolResult, error) {
	path, ok := req.Params.Arguments["path"].(string)
	if !ok {
		return mcp.NewToolResultError("path must be a string"), nil
	}
	pattern, ok := req.Params.Arguments["pattern"].(string)
	if !ok {
		return mcp.NewToolResultError("pattern must be a string"), nil
	}
	excludeMatcher := NewExcludeMatcher()
	excludePatterns, _ := req.Params.Arguments["excludePatterns"].([]interface{})
	for _, ep := range excludePatterns {
		epString := ep.(string)
		if epString == "" {
			continue
		}
		err := excludeMatcher.AddPattern(epString)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
	}
	validPath, err := validatePath(path, allowedDirs)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	// Return an error if the path is not a directory
	info, err := os.Stat(validPath)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if !info.IsDir() {
		return mcp.NewToolResultError("Path must be a directory"), nil
	}

	var results []string
	pattern = strings.ToLower(pattern)
	err = filepath.Walk(validPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		// Check relative path against exclude patterns
		if excludeMatcher.Match(validPath, filePath, info) {
			return nil
		}
		if strings.Contains(strings.ToLower(info.Name()), pattern) {
			results = append(results, filePath)
		}
		return nil
	})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if len(results) == 0 {
		return mcp.NewToolResultText("No matches found"), nil
	}
	return mcp.NewToolResultText(strings.Join(results, "\n")), nil
}
