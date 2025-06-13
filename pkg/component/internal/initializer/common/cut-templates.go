package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CutTemplate(file string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	documents := strings.Split(string(content), "\n---\n")
	var previousDocument string
	for _, doc := range documents {
		if strings.TrimSpace(doc) == "" {
			continue
		}

		lines := strings.Split(doc, "\n")
		var sourceLine string
		sourceIndex := 0

		for i := 0; i < 3 && i < len(lines); i++ {
			line := lines[i]
			if strings.HasPrefix(line, "# Source: ") {
				sourceLine = line
				sourceIndex = i
				break
			}
		}

		if sourceLine == "" {
			if previousDocument == "" {
				continue
			}
			sourceLine = previousDocument
		}

		originalPath := strings.TrimPrefix(sourceLine, "# Source: ")
		pathParts := strings.Split(originalPath, "/")
		filePath := filepath.Join(pathParts[2:]...)

		var contentBuilder strings.Builder
		contentBuilder.WriteString("---\n")

		lineCount := 0
		for i := sourceIndex; i < len(lines); i++ {
			line := lines[i]
			if !strings.HasPrefix(line, "---") && !strings.HasPrefix(line, "# Source: ") && line != "" {
				lineCount++
				contentBuilder.WriteString(line)
				contentBuilder.WriteString("\n")
			}
		}

		if lineCount == 0 {
			previousDocument = sourceLine
			continue
		}

		filePath = filepath.Join("base", filePath)
		dirPath := filepath.Dir(filePath)

		if err := os.MkdirAll(dirPath, 0755); err != nil {
			continue
		}

		finalPath := UniqueFilePath(filePath)

		if err := os.WriteFile(finalPath, []byte(contentBuilder.String()), 0644); err != nil {
			continue
		}
	}

	return nil
}
