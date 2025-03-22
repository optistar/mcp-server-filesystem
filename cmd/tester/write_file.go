package main

import (
	"github.com/mark3labs/mcp-go/mcp"
	"os"
	"path/filepath"
	"strings"
)

func writeFileTester(t T) {
	tempDir := t.TempDir()
	c := GetMCPClient(t.Context(), []string{tempDir})
	defer c.Close()

	// Create a subdirectory for testing
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create an existing file to test overwriting
	existingFilePath := filepath.Join(tempDir, "existing.txt")
	if err := os.WriteFile(existingFilePath, []byte("Original content"), 0644); err != nil {
		t.Fatalf("Failed to create existing test file: %v", err)
	}

	// Test cases
	testCases := []struct {
		name          string
		path          string
		content       string
		expectedError bool
		checkResult   func(string) bool
		verifyFile    func(string) bool
	}{
		{
			name:          "Create new file",
			path:          filepath.Join(tempDir, "newfile.txt"),
			content:       "New file content",
			expectedError: false,
			checkResult: func(result string) bool {
				return strings.Contains(result, "Successfully wrote")
			},
			verifyFile: func(path string) bool {
				content, err := os.ReadFile(path)
				return err == nil && string(content) == "New file content"
			},
		},
		{
			name:          "Overwrite existing file",
			path:          existingFilePath,
			content:       "Updated content",
			expectedError: false,
			checkResult: func(result string) bool {
				return strings.Contains(result, "Successfully wrote")
			},
			verifyFile: func(path string) bool {
				content, err := os.ReadFile(path)
				return err == nil && string(content) == "Updated content"
			},
		},
		{
			name:          "Write to subdirectory",
			path:          filepath.Join(subDir, "subfile.txt"),
			content:       "Content in subdirectory",
			expectedError: false,
			checkResult: func(result string) bool {
				return strings.Contains(result, "Successfully wrote")
			},
			verifyFile: func(path string) bool {
				content, err := os.ReadFile(path)
				return err == nil && string(content) == "Content in subdirectory"
			},
		},
		{
			name:          "Path outside allowed directories",
			path:          "/etc/passwd",
			content:       "This should fail",
			expectedError: true,
			checkResult: func(result string) bool {
				return strings.Contains(result, "access denied")
			},
			verifyFile: func(path string) bool {
				return true // Not checking file as it shouldn't be written
			},
		},
		{
			name:          "Write to non-existent parent directory",
			path:          filepath.Join(tempDir, "nonexistent/file.txt"),
			content:       "This should fail",
			expectedError: true,
			checkResult: func(result string) bool {
				return strings.Contains(result, "no such file or directory")
			},
			verifyFile: func(path string) bool {
				_, err := os.Stat(path)
				return os.IsNotExist(err) // File should not exist
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t T) {
			// Create request with path and content
			req := mcp.CallToolRequest{}
			req.Params.Name = "write_file"
			req.Params.Arguments = map[string]interface{}{
				"path":    tc.path,
				"content": tc.content,
			}

			// Call handler
			result, err := c.CallTool(t.Context(), req)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check result using helper function
			assertToolResult(t, result, tc.expectedError, tc.checkResult)

			// Verify file content
			if !tc.verifyFile(tc.path) {
				t.Errorf("File verification failed for path: %s", tc.path)
			}
		})
	}
}
