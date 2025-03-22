package top

import (
	"github.com/gobwas/glob"
	"os"
	"path/filepath"
	"strings"
)

type ExcludeMatcher interface {
	AddPattern(pattern string) error
	Match(dirPath string, filePath string, info os.FileInfo) bool
}

type excludeMatcher struct {
	matchers []patternMatcher
}

type patternMatcher struct {
	glob      glob.Glob
	isDirOnly bool
	anchored  bool
	prefix    string
	suffix    string
}

// NewExcludeMatcher creates a function that checks if a file should be excluded based on .gitignore patterns
func NewExcludeMatcher() ExcludeMatcher {
	return &excludeMatcher{}
}

func (e *excludeMatcher) AddPattern(pattern string) error {

	pattern = strings.TrimSpace(pattern)
	if pattern == "" || strings.HasPrefix(pattern, "#") {
		return nil // Skip empty lines and comments
	}

	isDirOnly := strings.HasSuffix(pattern, "/")
	anchored := strings.Contains(pattern, "/") && !strings.HasPrefix(pattern, "**/")

	// Clean up pattern for glob compilation
	compilePattern := pattern
	if isDirOnly {
		compilePattern = strings.TrimSuffix(pattern, "/")
	}

	// Handle special ** cases
	compilePattern = strings.ReplaceAll(compilePattern, "/**/", "[...]")
	if strings.HasPrefix(compilePattern, "**/") {
		compilePattern = strings.TrimPrefix(compilePattern, "**/")
	}
	if strings.HasSuffix(compilePattern, "/**") {
		compilePattern = strings.TrimSuffix(compilePattern, "/**")
	}

	g, err := glob.Compile(compilePattern, '/')
	if err != nil {
		return err
	}

	suffix := ""
	if strings.HasPrefix(pattern, "**/") {
		suffix = strings.TrimPrefix(pattern, "**/")
	}

	prefix := ""
	if strings.HasSuffix(pattern, "/**") {
		prefix = strings.TrimSuffix(pattern, "/**")
	}

	e.matchers = append(e.matchers, patternMatcher{
		glob:      g,
		isDirOnly: isDirOnly,
		anchored:  anchored,
		suffix:    suffix,
		prefix:    prefix,
	})
	return nil
}

func (e *excludeMatcher) Match(dirPath string, filePath string, info os.FileInfo) bool {
	// Convert to relative path
	relPath, err := filepath.Rel(dirPath, filePath)
	if err != nil {
		return false
	}
	// Convert to forward slashes for consistency
	relPath = filepath.ToSlash(relPath)

	for _, matcher := range e.matchers {
		// Check directory-only requirement
		if matcher.isDirOnly && !info.IsDir() {
			continue
		}

		// Handle anchored vs unanchored patterns
		if matcher.anchored {
			// For anchored patterns (containing /), match exact relative path
			if matcher.glob.Match(relPath) {
				return true
			}
		} else {
			// For unanchored patterns, check if any path segment matches
			segments := strings.Split(relPath, "/")
			for _, segment := range segments {
				if matcher.glob.Match(segment) {
					return true
				}
			}
		}

		// Handle **/prefix patterns by checking if path ends with pattern
		if matcher.suffix != "" && strings.HasSuffix(relPath, matcher.suffix) {
			return true
		}

		// Handle prefix/** patterns by checking if path starts with pattern
		if matcher.prefix != "" && strings.HasPrefix(relPath, matcher.prefix) {
			return true
		}
	}

	return false
}
