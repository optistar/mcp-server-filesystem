package top

import (
	"fmt"
	"github.com/pmezard/go-difflib/difflib"
	"os"
	"regexp"
	"strings"
)

// NormalizeLineEndings converts all Windows-style line endings to Unix-style
func normalizeLineEndings(text string) string {
	return strings.ReplaceAll(text, "\r\n", "\n")
}

// DetectLineEndingStyle determines the dominant line ending style in a text
// Returns "\r\n" for Windows style, "\n" for Unix style
func detectLineEndingStyle(text string) string {
	crlfCount := strings.Count(text, "\r\n")
	lfCount := strings.Count(text, "\n") - crlfCount // Subtract to avoid double counting

	if crlfCount > lfCount {
		return "\r\n"
	}
	return "\n"
}

func createUnifiedDiff(originalContent, newContent string, path string) (string, error) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(originalContent),
		B:        difflib.SplitLines(newContent),
		FromFile: path,
		ToFile:   path,
		Context:  3,
	}
	diffText, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return "", err
	}
	return diffText, nil
}

// Edit represents a single text replacement operation
type Edit struct {
	OldText string `json:"oldText"`
	NewText string `json:"newText"`
}

// ApplyFileEdits applies a series of edits to a file and returns a formatted diff
func applyFileEdits(originalPath, filePath string, edits []Edit, dryRun bool) (string, error) {
	// Read file content
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Detect original line ending style before normalization
	originalContent := string(contentBytes)
	lineEndingStyle := detectLineEndingStyle(originalContent)

	// Normalize line endings for processing
	content := normalizeLineEndings(originalContent)

	// Apply edits sequentially
	modifiedContent := content
	for _, edit := range edits {
		normalizedOld := normalizeLineEndings(edit.OldText)
		normalizedNew := normalizeLineEndings(edit.NewText)

		// If exact match exists, use it
		if strings.Contains(modifiedContent, normalizedOld) {
			modifiedContent = strings.Replace(modifiedContent, normalizedOld, normalizedNew, 1)
			continue
		}

		// Otherwise, try line-by-line matching with flexibility for whitespace
		oldLines := strings.Split(normalizedOld, "\n")
		contentLines := strings.Split(modifiedContent, "\n")
		matchFound := false

		for i := 0; i <= len(contentLines)-len(oldLines); i++ {
			potentialMatch := contentLines[i : i+len(oldLines)]

			// Compare lines with normalized whitespace
			isMatch := true
			for j, oldLine := range oldLines {
				if strings.TrimSpace(oldLine) != strings.TrimSpace(potentialMatch[j]) {
					isMatch = false
					break
				}
			}

			if isMatch {
				// Preserve original indentation of first line
				re := regexp.MustCompile(`^\s*`)
				originalIndent := re.FindString(contentLines[i])

				newLines := strings.Split(normalizedNew, "\n")
				for j := range newLines {
					if j == 0 {
						newLines[j] = originalIndent + strings.TrimLeft(newLines[j], " \t")
					} else {
						// For subsequent lines, try to preserve relative indentation
						var oldIndent string
						if j < len(oldLines) {
							oldIndent = re.FindString(oldLines[j])
						}
						newIndent := re.FindString(newLines[j])
						if oldIndent != "" && newIndent != "" {
							relativeIndent := len(newIndent) - len(oldIndent)
							if relativeIndent > 0 {
								newLines[j] = originalIndent + strings.Repeat(" ", relativeIndent) +
									strings.TrimLeft(newLines[j], " \t")
							} else {
								newLines[j] = originalIndent + strings.TrimLeft(newLines[j], " \t")
							}
						}
					}
				}

				// Replace the matching section
				contentLines = append(contentLines[:i],
					append(newLines, contentLines[i+len(oldLines):]...)...)
				modifiedContent = strings.Join(contentLines, "\n")
				matchFound = true
				break
			}
		}

		if !matchFound {
			return "", fmt.Errorf("could not find exact match for edit:\n%s", edit.OldText)
		}
	}

	// Create unified diff
	diff, err := createUnifiedDiff(content, modifiedContent, originalPath)
	if err != nil {
		return "", fmt.Errorf("failed to create diff: %w", err)
	}

	// Format diff with appropriate number of backticks
	numBackticks := 3
	for strings.Contains(diff, strings.Repeat("`", numBackticks)) {
		numBackticks++
	}
	formattedDiff := fmt.Sprintf("%sdiff\n%s%s\n\n",
		strings.Repeat("`", numBackticks),
		diff,
		strings.Repeat("`", numBackticks))

	if !dryRun {
		// Convert back to original line ending style before writing
		finalContent := modifiedContent
		if lineEndingStyle == "\r\n" {
			finalContent = strings.ReplaceAll(modifiedContent, "\n", "\r\n")
		}

		err = os.WriteFile(filePath, []byte(finalContent), 0644)
		if err != nil {
			return "", err
		}
	}

	return formattedDiff, nil
}
