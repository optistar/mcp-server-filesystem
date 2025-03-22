package top

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
)

func DefineEditFileTool() mcp.Tool {
	return mcp.NewTool("edit_file",
		mcp.WithDescription(
			"Make line-based edits to a text file. Each edit replaces exact line sequences "+
				"with new content. Returns a git-style diff showing the changes made. "+
				"Only works within allowed directories."),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to the file")),
		mcp.WithArray("edits",
			mcp.Required(),
			mcp.Description("Array of edit operations"),
			mcp.Items(map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"oldText": map[string]interface{}{
						"type":        "string",
						"description": "Text to replace",
					},
					"newText": map[string]interface{}{
						"type":        "string",
						"description": "Replacement text",
					},
				},
				"required": []string{"oldText", "newText"},
			}),
		),
		mcp.WithBoolean("dryRun",
			mcp.Description("Preview changes using git-style diff format"),
			mcp.DefaultBool(false),
		),
	)
}

func editsFromJSON(edits interface{}) ([]Edit, error) {
	var result []Edit

	// Handle both possible types
	switch editsTyped := edits.(type) {
	case []map[string]interface{}:
		// Handle direct []map[string]interface{} type
		for _, editMap := range editsTyped {
			oldText, ok := editMap["oldText"].(string)
			if !ok {
				return nil, fmt.Errorf("oldText must be a string")
			}

			newText, ok := editMap["newText"].(string)
			if !ok {
				return nil, fmt.Errorf("newText must be a string")
			}

			result = append(result, Edit{OldText: oldText, NewText: newText})
		}
	case []interface{}:
		// Handle []interface{} type
		for _, edit := range editsTyped {
			editMap, ok := edit.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("edit must be a map")
			}

			oldText, ok := editMap["oldText"].(string)
			if !ok {
				return nil, fmt.Errorf("oldText must be a string")
			}

			newText, ok := editMap["newText"].(string)
			if !ok {
				return nil, fmt.Errorf("newText must be a string")
			}

			result = append(result, Edit{OldText: oldText, NewText: newText})
		}
	default:
		return nil, fmt.Errorf("edits must be an array")
	}

	return result, nil
}

func EditFileHandler(_ context.Context, req mcp.CallToolRequest, allowedDirs []string) (*mcp.CallToolResult, error) {
	path, ok := req.Params.Arguments["path"].(string)
	if !ok {
		return mcp.NewToolResultError("path must be a string"), nil
	}

	// Get the edits from the request arguments
	editsRaw, ok := req.Params.Arguments["edits"]
	if !ok {
		return mcp.NewToolResultError("edits parameter is required"), nil
	}

	// Handle both possible types: []interface{} and []map[string]interface{}
	var edits []Edit
	var err error

	// Try to convert directly to []map[string]interface{} first
	if editsArray, ok := editsRaw.([]map[string]interface{}); ok {
		edits, err = editsFromJSON(editsArray)
	} else if editsArray, ok := editsRaw.([]interface{}); ok {
		// Fall back to the original approach
		edits, err = editsFromJSON(editsArray)
	} else {
		return mcp.NewToolResultError("edits must be an array"), nil
	}

	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	dryRun, _ := req.Params.Arguments["dryRun"].(bool)

	validPath, err := validatePath(path, allowedDirs)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	diffText, err := applyFileEdits(path, validPath, edits, dryRun)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(diffText), nil
}
