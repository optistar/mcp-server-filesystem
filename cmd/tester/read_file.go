package main

import (
	"github.com/mark3labs/mcp-go/mcp"
	"os"
	"path/filepath"
)

func readFileTester(t T) {
	tempDir := t.TempDir()
	c := GetMCPClient(t.Context(), []string{tempDir})
	defer c.Close()

	// Create a subdirectory for testing path traversal
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create a test file with content
	validFilePath := filepath.Join(tempDir, "test.txt")
	testContent := "This is test content"
	if err := os.WriteFile(validFilePath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a test file in subdirectory
	subDirFilePath := filepath.Join(subDir, "subfile.txt")
	if err := os.WriteFile(subDirFilePath, []byte("Content in subdirectory"), 0644); err != nil {
		t.Fatalf("Failed to create test file in subdirectory: %v", err)
	}

	// Test cases
	tests := []struct {
		name          string
		path          string
		expectedError bool
		expectedText  string
	}{
		{
			name:          "Valid file",
			path:          validFilePath,
			expectedError: false,
			expectedText:  testContent,
		},
		{
			name:          "Non-existent file",
			path:          filepath.Join(tempDir, "nonexistent.txt"),
			expectedError: true,
			expectedText:  "", // Error expected
		},
		{
			name:          "Path outside allowed directories",
			path:          "/etc/passwd", // This should be denied
			expectedError: true,
			expectedText:  "", // Error expected
		},
		{
			name:          "Path traversal attempt",
			path:          filepath.Join(subDir, "../test.txt"),
			expectedError: false, // This should work as it resolves to a valid path
			expectedText:  testContent,
		},
		{
			name:          "Path traversal outside allowed dirs",
			path:          filepath.Join(tempDir, "../../../etc/passwd"),
			expectedError: true,
			expectedText:  "", // Error expected
		},
		{
			name:          "Subdirectory file",
			path:          subDirFilePath,
			expectedError: false,
			expectedText:  "Content in subdirectory",
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t T) {
			// Create request with path
			req := mcp.CallToolRequest{}
			req.Params.Name = "read_file"
			req.Params.Arguments = map[string]interface{}{
				"path": tc.path,
			}

			// Call handler
			result, err := c.CallTool(t.Context(), req)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check result using helper function
			assertToolResult(t, result, tc.expectedError, expectExactText(tc.expectedText))
		})
	}
}
