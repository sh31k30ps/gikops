package helm

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sh31k30ps/gikopsctl/pkg/config/component"
)

func handleUpload(upload component.HelmHookUpload) error {
	resp, err := http.Get(upload.URL)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	outPath := filepath.Join("base", upload.Name)
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", filepath.Dir(outPath), err)
	}
	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func handleResolve(pattern string) error {
	files, err := filepath.Glob(filepath.Join("base", pattern))
	if err != nil {
		return fmt.Errorf("failed to glob files: %w", err)
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		// Find URL in content and download
		// This is a simplified version - you might want to improve the URL detection
		if url := findURL(string(content)); url != "" {
			resp, err := http.Get(url)
			if err != nil {
				return fmt.Errorf("failed to download from URL: %w", err)
			}
			defer resp.Body.Close()

			additionalContent, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response: %w", err)
			}

			// Append the downloaded content
			content = append(content, additionalContent...)
			if err := os.WriteFile(file, content, 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
		}
	}

	return nil
}

func handleRename(rename component.HelmHookRename) error {
	oldPath := filepath.Join("base", rename.Original)
	newPath := filepath.Join("base", rename.Renamed)
	return os.Rename(oldPath, newPath)
}

func handleConcat(concat component.HelmHookConcat) error {
	if len(concat.Includes) == 0 {
		concat.Includes = []string{"*.yaml"}
	}
	var files []string
	for _, pattern := range concat.Includes {
		detectedfiles, err := filepath.Glob(filepath.Join("base", concat.Folder, pattern))
		if err != nil {
			return fmt.Errorf("failed to glob files with pattern %s: %w", concat.Includes, err)
		}
		if len(files) == 0 {
			files = detectedfiles
			continue
		}
		files = append(files, detectedfiles...)
	}
	outputFile := filepath.Join("base", concat.Output)
	outputDir := filepath.Dir(filepath.Join("base", outputFile))
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputFile, err)
	}
	defer output.Close()

	// Concatenate all files
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", file, err)
		}

		if _, err := output.Write(content); err != nil {
			return fmt.Errorf("failed to write to output file: %w", err)
		}

		if err := os.Remove(file); err != nil {
			return fmt.Errorf("failed to remove file '%s': %w", file, err)
		}
	}

	return nil
}

func findURL(content string) string {
	// Simple URL detection - you might want to use a more robust solution
	if idx := strings.Index(content, "http"); idx != -1 {
		end := strings.IndexAny(content[idx:], " \n\t")
		if end == -1 {
			return content[idx:]
		}
		return content[idx : idx+end]
	}
	return ""
}
