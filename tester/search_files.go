package tester

import (
	"github.com/mark3labs/mcp-go/mcp"
	"os"
	"path/filepath"
	"strings"
)

func TestSearchFiles(t T, f MCPClientFactory) {
	tempDir := t.TempDir()
	_, c := f(t.Context(), []string{tempDir})
	defer c.Close()

	// Create a nested directory structure for testing
	dirs := []string{
		filepath.Join(tempDir, "dir1"),
		filepath.Join(tempDir, "dir1/subdir1"),
		filepath.Join(tempDir, "dir2"),
		filepath.Join(tempDir, "dir2/subdir2"),
		filepath.Join(tempDir, "dir3"),
		filepath.Join(tempDir, "dir3/subdir3"),
		filepath.Join(tempDir, "emptydir"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create test files with various names for pattern matching
	testFiles := map[string]string{
		filepath.Join(tempDir, "file1.txt"):               "Content of file 1",
		filepath.Join(tempDir, "file2.txt"):               "Content of file 2",
		filepath.Join(tempDir, "document.doc"):            "Document content",
		filepath.Join(tempDir, "dir1/file1.txt"):          "Dir1 file1",
		filepath.Join(tempDir, "dir1/test.txt"):           "Dir1 test file",
		filepath.Join(tempDir, "dir1/subdir1/test.txt"):   "Subdir1 test file",
		filepath.Join(tempDir, "dir2/file2.txt"):          "Dir2 file2",
		filepath.Join(tempDir, "dir2/test.doc"):           "Dir2 test doc",
		filepath.Join(tempDir, "dir2/subdir2/test.txt"):   "Subdir2 test file",
		filepath.Join(tempDir, "dir3/file3.txt"):          "Dir3 file3",
		filepath.Join(tempDir, "dir3/subdir3/hidden.txt"): "Hidden file",
	}

	for path, content := range testFiles {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", path, err)
		}
	}

	// Test cases
	testCases := []struct {
		name            string
		path            string
		pattern         string
		excludePatterns []interface{}
		expectedError   bool
		checkResult     func(string) bool
	}{
		{
			name:          "Search for .txt files",
			path:          tempDir,
			pattern:       ".txt",
			expectedError: false,
			checkResult: func(result string) bool {
				// Should find all .txt files
				return strings.Count(result, ".txt") >= 8 && // At least 8 .txt files
					strings.Contains(result, "file1.txt") &&
					strings.Contains(result, "test.txt")
			},
		},
		{
			name:          "Search for specific file name",
			path:          tempDir,
			pattern:       "file1",
			expectedError: false,
			checkResult: func(result string) bool {
				return strings.Contains(result, "file1.txt") &&
					!strings.Contains(result, "file2.txt") &&
					strings.Count(result, "file1") >= 2 // Root and dir1
			},
		},
		{
			name:          "Search in specific subdirectory",
			path:          filepath.Join(tempDir, "dir2"),
			pattern:       "test",
			expectedError: false,
			checkResult: func(result string) bool {
				return strings.Contains(result, "test.doc") &&
					strings.Contains(result, "test.txt") &&
					!strings.Contains(result, "dir1") // Should not search outside dir2
			},
		},
		{
			name:          "Search with no matches",
			path:          tempDir,
			pattern:       "nonexistent",
			expectedError: false,
			checkResult: func(result string) bool {
				return strings.Contains(result, "No matches found")
			},
		},
		{
			name:          "Search in empty directory",
			path:          filepath.Join(tempDir, "emptydir"),
			pattern:       ".",
			expectedError: false,
			checkResult: func(result string) bool {
				return strings.Contains(result, "No matches found")
			},
		},
		{
			name:            "Search with exclude pattern",
			path:            tempDir,
			pattern:         "test",
			excludePatterns: []interface{}{"dir1/**"},
			expectedError:   false,
			checkResult: func(result string) bool {
				return strings.Contains(result, "dir2") &&
					!strings.Contains(result, "dir1/test.txt") &&
					!strings.Contains(result, "dir1/subdir1/test.txt")
			},
		},
		{
			name:            "Search with multiple exclude patterns",
			path:            tempDir,
			pattern:         ".txt",
			excludePatterns: []interface{}{"dir1/**", "dir2/**"},
			expectedError:   false,
			checkResult: func(result string) bool {
				return strings.Contains(result, "file1.txt") &&
					strings.Contains(result, "file2.txt") &&
					strings.Contains(result, "dir3") &&
					!strings.Contains(result, "dir1/") &&
					!strings.Contains(result, "dir2/")
			},
		},
		{
			name:            "Search with simple filename exclude",
			path:            tempDir,
			pattern:         "file",
			excludePatterns: []interface{}{"file1.txt"},
			expectedError:   false,
			checkResult: func(result string) bool {
				return !strings.Contains(result, "/file1.txt") &&
					strings.Contains(result, "file2.txt") &&
					strings.Contains(result, "file3.txt")
			},
		},
		{
			name:          "Path outside allowed directories",
			path:          "/etc",
			pattern:       "passwd",
			expectedError: true,
			checkResult: func(result string) bool {
				return strings.Contains(result, "access denied")
			},
		},
		{
			name:          "Path traversal attempt",
			path:          filepath.Join(tempDir, "dir1/../dir2"),
			pattern:       "test",
			expectedError: false, // Should resolve to valid path
			checkResult: func(result string) bool {
				return strings.Contains(result, "test.doc") &&
					strings.Contains(result, "test.txt")
			},
		},
		{
			name:          "Non-existent directory",
			path:          filepath.Join(tempDir, "nonexistent"),
			pattern:       "test",
			expectedError: true,
			checkResult: func(result string) bool {
				return strings.Contains(result, "no such file or directory")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t T) {
			// Create request with path and pattern
			req := mcp.CallToolRequest{}
			req.Params.Name = "search_files"
			req.Params.Arguments = map[string]interface{}{
				"path":    tc.path,
				"pattern": tc.pattern,
			}

			// Add exclude patterns if provided
			if tc.excludePatterns != nil {
				req.Params.Arguments["excludePatterns"] = tc.excludePatterns
			}

			// Call handler
			result, err := c.CallTool(t.Context(), req)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check result using helper function
			assertToolResult(t, result, tc.expectedError, tc.checkResult)
		})
	}
}
