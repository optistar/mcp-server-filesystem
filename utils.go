package top

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func isSymlink(info os.FileInfo) bool {
	return info.Mode()&os.ModeSymlink == os.ModeSymlink
}

func isInAllowedDirectories(cleanPath string, allowedDirectories []string) (bool, string) {
	for _, dir := range allowedDirectories {
		if strings.HasPrefix(cleanPath, dir) {
			info, err := os.Lstat(cleanPath)
			if err == nil {
				// Path exists - validate it directly
				if isSymlink(info) {
					linkTarget, err := os.Readlink(cleanPath)
					if err != nil {
						return false, ""
					}
					return true, linkTarget
				}
				return true, ""
			}
			// For non-existent paths, we need to find the closest existing parent
			currentPath := cleanPath
			for {
				parentPath := filepath.Dir(currentPath)

				// If we've reached the root or gone outside allowed dir, stop
				if parentPath == currentPath || !strings.HasPrefix(parentPath, dir) {
					return false, ""
				}

				// Check if this parent exists
				parentInfo, err := os.Lstat(parentPath)
				if err == nil {
					// Path exists - validate it directly
					if isSymlink(parentInfo) {
						parentTarget, err := os.Readlink(parentPath)
						if err != nil {
							return false, ""
						}
						return true, filepath.Join(parentTarget, strings.TrimPrefix(cleanPath, parentPath))
					}
					return true, ""
				}

				// Move up to the next parent
				currentPath = parentPath
			}
		}
	}
	return false, ""
}

func validatePath(requestedPath string, allowedDirectories []string) (string, error) {
	absPath, err := filepath.Abs(ExpandHome(requestedPath))
	if err != nil {
		return "", fmt.Errorf("invalid path: %v", err)
	}
	cleanPath := filepath.Clean(absPath)

	// Set for loop detection
	visited := map[string]struct{}{}
	// Check that symlinks can be resolved to a real file  while staying inside allowed directories
	tempPath := cleanPath
	for {
		visited[tempPath] = struct{}{}
		ok, target := isInAllowedDirectories(tempPath, allowedDirectories)
		if !ok {
			return "", fmt.Errorf("access denied - path outside allowed directories: %s", absPath)
		}
		if target == "" {
			// No symlink - we're done
			return cleanPath, nil
		}
		// Follow the symlink
		if _, ok := visited[target]; ok {
			return "", fmt.Errorf("access denied - symlink loop detected: %s", absPath)
		}
		tempPath = target
	}
}

// Utility functions
func ExpandHome(path string) string {
	if strings.HasPrefix(path, "~/") || path == "~" {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[1:])
	}
	return path
}
