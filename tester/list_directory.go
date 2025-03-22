package tester

import (
	"github.com/mark3labs/mcp-go/mcp"
	"os"
	"path/filepath"
	"strings"
)

func TestListDirectory(t T, f MCPClientFactory) {
	tempDir := t.TempDir()
	_, c := f(t.Context(), []string{tempDir})
	defer c.Close()

	// Create subdirectories for testing
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	emptyDir := filepath.Join(tempDir, "emptydir")
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatalf("Failed to create empty directory: %v", err)
	}

	// Create test files
	testFiles := map[string]string{
		filepath.Join(tempDir, "file1.txt"):  "Content of file 1",
		filepath.Join(tempDir, "file2.txt"):  "Content of file 2",
		filepath.Join(subDir, "subfile.txt"): "Content in subdirectory",
	}

	for path, content := range testFiles {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", path, err)
		}
	}

	// Test cases
	testCases := []struct {
		name          string
		path          string
		expectedError bool
		checkContent  func(string) bool
	}{
		{
			name:          "List directory with files and subdirectories",
			path:          tempDir,
			expectedError: false,
			checkContent: func(content string) bool {
				return strings.Contains(content, "[FILE] file1.txt") &&
					strings.Contains(content, "[FILE] file2.txt") &&
					strings.Contains(content, "[DIR] subdir") &&
					strings.Contains(content, "[DIR] emptydir")
			},
		},
		{
			name:          "List subdirectory",
			path:          subDir,
			expectedError: false,
			checkContent: func(content string) bool {
				return strings.Contains(content, "[FILE] subfile.txt")
			},
		},
		{
			name:          "List empty directory",
			path:          emptyDir,
			expectedError: false,
			checkContent: func(content string) bool {
				return content == "Empty directory" // Empty directory should return empty content
			},
		},
		{
			name:          "Non-existent directory",
			path:          filepath.Join(tempDir, "nonexistent"),
			expectedError: true,
			checkContent: func(content string) bool {
				return strings.Contains(content, "no such file or directory")
			},
		},
		{
			name:          "Path is a file, not a directory",
			path:          filepath.Join(tempDir, "file1.txt"),
			expectedError: true,
			checkContent: func(content string) bool {
				return strings.Contains(content, "not a directory")
			},
		},
		{
			name:          "Path outside allowed directories",
			path:          "/etc",
			expectedError: true,
			checkContent: func(content string) bool {
				return strings.Contains(content, "access denied")
			},
		},
		{
			name:          "Path traversal attempt",
			path:          filepath.Join(subDir, ".."),
			expectedError: false, // Should resolve to tempDir and work
			checkContent: func(content string) bool {
				return strings.Contains(content, "[FILE] file1.txt") &&
					strings.Contains(content, "[FILE] file2.txt")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t T) {
			// Create request with path
			req := mcp.CallToolRequest{}
			req.Params.Name = "list_directory"
			req.Params.Arguments = map[string]interface{}{
				"path": tc.path,
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
