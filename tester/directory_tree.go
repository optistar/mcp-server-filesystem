package tester

import (
	"encoding/json"
	"github.com/mark3labs/mcp-go/mcp"
	"os"
	"path/filepath"
	"strings"
)

func TestDirectoryTree(t T, f MCPClientFactory) {
	tempDir := t.TempDir()
	_, c := f(t.Context(), []string{tempDir})
	defer c.Close()

	// Create a nested directory structure for testing
	dirs := []string{
		filepath.Join(tempDir, "dir1"),
		filepath.Join(tempDir, "dir1/subdir1"),
		filepath.Join(tempDir, "dir1/subdir2"),
		filepath.Join(tempDir, "dir2"),
		filepath.Join(tempDir, "dir2/subdir1"),
		filepath.Join(tempDir, "emptydir"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create test files
	testFiles := map[string]string{
		filepath.Join(tempDir, "file1.txt"):              "Root file",
		filepath.Join(tempDir, "dir1/file1.txt"):         "Dir1 file",
		filepath.Join(tempDir, "dir1/subdir1/file1.txt"): "Subdir1 file",
		filepath.Join(tempDir, "dir1/subdir2/file2.txt"): "Subdir2 file",
		filepath.Join(tempDir, "dir2/file3.txt"):         "Dir2 file",
		filepath.Join(tempDir, "dir2/subdir1/file4.txt"): "Dir2/subdir1 file",
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
			name:          "Root directory tree",
			path:          tempDir,
			expectedError: false,
			checkContent: func(content string) bool {
				// Check for JSON structure with all directories and files
				return strings.Contains(content, `"name": "file1.txt"`) &&
					strings.Contains(content, `"name": "dir1"`) &&
					strings.Contains(content, `"name": "dir2"`) &&
					strings.Contains(content, `"name": "emptydir"`) &&
					strings.Contains(content, `"type": "directory"`) &&
					strings.Contains(content, `"type": "file"`) &&
					strings.Contains(content, `"children"`)
			},
		},
		{
			name:          "Subdirectory tree",
			path:          filepath.Join(tempDir, "dir1"),
			expectedError: false,
			checkContent: func(content string) bool {
				return strings.Contains(content, `"name": "file1.txt"`) &&
					strings.Contains(content, `"name": "subdir1"`) &&
					strings.Contains(content, `"name": "subdir2"`) &&
					!strings.Contains(content, `"name": "dir2"`) // Should not include dir2
			},
		},
		{
			name:          "Empty directory tree",
			path:          filepath.Join(tempDir, "emptydir"),
			expectedError: false,
			checkContent: func(content string) bool {
				return content == "[]"
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
			path:          filepath.Join(tempDir, "dir1/.."),
			expectedError: false, // Should resolve to tempDir and work
			checkContent: func(content string) bool {
				return strings.Contains(content, `"name": "file1.txt"`) &&
					strings.Contains(content, `"name": "dir1"`) &&
					strings.Contains(content, `"name": "dir2"`)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t T) {
			// Create request with path
			req := mcp.CallToolRequest{}
			req.Params.Name = "directory_tree"
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

			// For valid results, verify JSON can be parsed
			if !tc.expectedError && result != nil && !result.IsError {
				textContent, ok := result.Content[0].(mcp.TextContent)
				if !ok {
					t.Fatalf("Expected text content but got: %v", result.Content)
				}

				// Try to parse the JSON
				var treeData interface{}
				if err := json.Unmarshal([]byte(textContent.Text), &treeData); err != nil {
					t.Errorf("Failed to parse JSON response: %v\nJSON: %s", err, textContent.Text)
				}
			}
		})
	}
}
