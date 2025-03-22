package main

import (
	"github.com/mark3labs/mcp-go/mcp"
	"os"
	"path/filepath"
)

func editFileTester(t T) {
	tempDir := t.TempDir()
	c := GetMCPClient(t.Context(), []string{tempDir})
	defer c.Close()

	// Create a test file with mixed line endings
	testFilePath := filepath.Join(tempDir, "test.txt")
	originalContent := "Line 1\r\nLine 2\r\nLine 3\nLine 4\r\n"
	err := os.WriteFile(testFilePath, []byte(originalContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test 1: Valid edit request
	edits := []map[string]interface{}{
		{
			"oldText": "Line 2",
			"newText": "Modified Line 2",
		},
	}

	req := mcp.CallToolRequest{}
	req.Params.Name = "edit_file"
	req.Params.Arguments = map[string]interface{}{
		"path":   testFilePath,
		"edits":  edits,
		"dryRun": false,
	}

	result, err := c.CallTool(t.Context(), req)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if result.IsError {
		t.Errorf("Expected success but got error: %s", result.Content)
	}

	// Verify file was changed correctly
	contentAfterEdit, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read file after edit: %v", err)
	}

	expectedContent := "Line 1\r\nModified Line 2\r\nLine 3\r\nLine 4\r\n"
	if string(contentAfterEdit) != expectedContent {
		t.Errorf("File content doesn't match expected.\nGot: %q\nWant: %q",
			string(contentAfterEdit), expectedContent)
	}

	// Test 2: Dry run should not modify the file
	edits = []map[string]interface{}{
		{
			"oldText": "Line 3",
			"newText": "Modified Line 3",
		},
	}

	req.Params.Arguments = map[string]interface{}{
		"path":   testFilePath,
		"edits":  edits,
		"dryRun": true,
	}

	result, err = c.CallTool(t.Context(), req)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if result.IsError {
		t.Errorf("Expected success but got error: %s", result.Content)
	}

	// Verify file was NOT changed (should still have content from first edit)
	contentAfterDryRun, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read file after dry run: %v", err)
	}

	if string(contentAfterDryRun) != expectedContent {
		t.Errorf("File content was modified during dry run.\nGot: %q\nWant: %q",
			string(contentAfterDryRun), expectedContent)
	}

	// Test 3: Invalid path
	req = mcp.CallToolRequest{}
	req.Params.Name = "edit_file"
	req.Params.Arguments = map[string]interface{}{
		"path":   filepath.Join(tempDir, "nonexistent", "file.txt"),
		"edits":  edits,
		"dryRun": false,
	}

	result, err = c.CallTool(t.Context(), req)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if !result.IsError {
		t.Errorf("Expected error for invalid path but got success: %s", result.Content)
	}

	// Test 4: Invalid edits format
	req.Params.Arguments = map[string]interface{}{
		"path":   testFilePath,
		"edits":  "not an array",
		"dryRun": false,
	}

	result, err = c.CallTool(t.Context(), req)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if !result.IsError {
		t.Errorf("Expected error for invalid edits format but got success: %s", result.Content)
	}

	// Test 5: Missing oldText in edit
	edits = []map[string]interface{}{
		{
			"newText": "Only new text provided",
		},
	}

	req.Params.Arguments = map[string]interface{}{
		"path":   testFilePath,
		"edits":  edits,
		"dryRun": false,
	}

	result, err = c.CallTool(t.Context(), req)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if !result.IsError {
		t.Errorf("Expected error for missing oldText but got success: %s", result.Content)
	}
}
