package main

import (
	"github.com/mark3labs/mcp-go/mcp"
	"os"
	"path/filepath"
	"strings"
)

func readMultipleFilesTester(t T) {
	tempDir := t.TempDir()
	c := GetMCPClient(t.Context(), []string{tempDir})
	defer c.Close()

	// Create a subdirectory for testing
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create test files
	testFiles := map[string]string{
		filepath.Join(tempDir, "file1.txt"):      "Content of file 1",
		filepath.Join(tempDir, "file2.txt"):      "Content of file 2",
		filepath.Join(subDir, "subdir_file.txt"): "Content in subdirectory",
	}

	for path, content := range testFiles {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", path, err)
		}
	}

	// Test cases
	testCases := []struct {
		name          string
		paths         []interface{}
		expectedError bool
		checkContent  func(string) bool
	}{
		{
			name:          "Valid multiple files",
			paths:         []interface{}{filepath.Join(tempDir, "file1.txt"), filepath.Join(tempDir, "file2.txt")},
			expectedError: false,
			checkContent: func(content string) bool {
				return strings.Contains(content, "Content of file 1") &&
					strings.Contains(content, "Content of file 2")
			},
		},
		{
			name:          "Mix of valid and invalid paths",
			paths:         []interface{}{filepath.Join(tempDir, "file1.txt"), filepath.Join(tempDir, "nonexistent.txt")},
			expectedError: false,
			checkContent: func(content string) bool {
				return strings.Contains(content, "Content of file 1") &&
					strings.Contains(content, "Error")
			},
		},
		{
			name:          "Path outside allowed directories",
			paths:         []interface{}{filepath.Join(tempDir, "file1.txt"), "/etc/passwd"},
			expectedError: false,
			checkContent: func(content string) bool {
				return strings.Contains(content, "Content of file 1") &&
					strings.Contains(content, "access denied")
			},
		},
		{
			name:          "Non-string path",
			paths:         []interface{}{filepath.Join(tempDir, "file1.txt"), 123},
			expectedError: false,
			checkContent: func(content string) bool {
				return strings.Contains(content, "Content of file 1") &&
					strings.Contains(content, "Error - must be a string")
			},
		},
		{
			name:          "Subdirectory file",
			paths:         []interface{}{filepath.Join(subDir, "subdir_file.txt")},
			expectedError: false,
			checkContent: func(content string) bool {
				return strings.Contains(content, "Content in subdirectory")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t T) {
			// Create request with paths
			req := mcp.CallToolRequest{}
			req.Params.Name = "read_multiple_files"
			req.Params.Arguments = map[string]interface{}{
				"paths": tc.paths,
			}

			// Call handler
			result, err := c.CallTool(t.Context(), req)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check result using helper function
			assertToolResult(t, result, tc.expectedError, tc.checkContent)
		})
	}
}
