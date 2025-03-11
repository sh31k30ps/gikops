package assets

import (
	"embed"
	"fmt"
	"path"
	"strings"
)

//go:embed kind.yaml gitignore overrides/** components/**
var Assets embed.FS

// GetKindConfig returns the kind configuration file content
func GetKindConfig() ([]byte, error) {
	return Assets.ReadFile("kind.yaml")
}

// GetGitignore returns the gitignore file content
func GetGitignore() ([]byte, error) {
	return Assets.ReadFile("gitignore")
}

// getFilesFromDir returns all files from a directory in the embedded assets
func getFilesFromDir(dir string) ([]string, error) {
	entries, err := Assets.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, path.Join(dir, entry.Name()))
			continue
		}
		subFiles, err := getFilesFromDir(path.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		files = append(files, subFiles...)
	}
	return files, nil
}

// GetSubdirectories returns all subdirectories from a directory in the embedded assets
func GetSubdirectories(dir string) ([]string, error) {
	entries, err := Assets.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, path.Join(dir, entry.Name()))
		}
	}
	return dirs, nil
}

// GetFilesFromSubdirectory returns all files from a specific subdirectory
func GetFilesFromSubdirectory(subdir string) ([]string, error) {
	if subdir == "" {
		return nil, fmt.Errorf("subdirectory path cannot be empty")
	}
	return getFilesFromDir(subdir)
}

// GetComponentFiles returns all files from the components directory
func GetComponentsFiles() ([]string, error) {
	return getFilesFromDir("components")
}

// GetComponentFiles returns all files from the components directory
func GetComponentFiles(component string) ([]string, error) {
	return getFilesFromDir(path.Join("components", component))
}

// GetOverrideFiles returns all files from the overrides directory
func GetOverrideFiles() ([]string, error) {
	return getFilesFromDir("overrides")
}

// GetFile returns the content of a specific file from the embedded assets
func GetFile(filepath string) ([]byte, error) {
	if filepath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}
	return Assets.ReadFile(filepath)
}

// GetFilesByExtension returns all files with a specific extension from a directory
func GetFilesByExtension(dir, ext string) ([]string, error) {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	files, err := getFilesFromDir(dir)
	if err != nil {
		return nil, err
	}

	var filtered []string
	for _, file := range files {
		if strings.HasSuffix(file, ext) {
			filtered = append(filtered, file)
		}
	}
	return filtered, nil
}
