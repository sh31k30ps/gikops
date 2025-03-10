package assets

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestGetKindConfig(t *testing.T) {
	content, err := GetKindConfig()
	if err != nil {
		t.Errorf("GetKindConfig() error = %v", err)
		return
	}
	if len(content) == 0 {
		t.Error("GetKindConfig() returned empty content")
	}
}

func TestGetComponentFiles(t *testing.T) {
	files, err := GetComponentFiles()
	if err != nil {
		t.Errorf("GetComponentFiles() error = %v", err)
		return
	}

	// Verify all paths start with components/
	for _, file := range files {
		if !strings.HasPrefix(file, "components/") {
			t.Errorf("GetComponentFiles() returned file with invalid prefix: %s", file)
		}
	}
}

func TestGetOverrideFiles(t *testing.T) {
	files, err := GetOverrideFiles()
	if err != nil {
		t.Errorf("GetOverrideFiles() error = %v", err)
		return
	}

	// Verify all paths start with overrides/
	for _, file := range files {
		if !strings.HasPrefix(file, "overrides/") {
			t.Errorf("GetOverrideFiles() returned file with invalid prefix: %s", file)
		}
	}
}

func TestGetSubdirectories(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		wantErr bool
	}{
		{
			name:    "Valid directory",
			dir:     "components",
			wantErr: false,
		},
		{
			name:    "Invalid directory",
			dir:     "nonexistent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dirs, err := GetSubdirectories(tt.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSubdirectories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				for _, dir := range dirs {
					if !strings.HasPrefix(dir, tt.dir+"/") {
						t.Errorf("GetSubdirectories() returned dir with invalid prefix: %s", dir)
					}
				}
			}
		})
	}
}

func TestGetFilesFromSubdirectory(t *testing.T) {
	tests := []struct {
		name    string
		subdir  string
		wantErr bool
	}{
		{
			name:    "Empty subdirectory",
			subdir:  "",
			wantErr: true,
		},
		{
			name:    "Invalid subdirectory",
			subdir:  "nonexistent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := GetFilesFromSubdirectory(tt.subdir)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFilesFromSubdirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(files) > 0 {
				for _, file := range files {
					if !strings.HasPrefix(file, tt.subdir+"/") {
						t.Errorf("GetFilesFromSubdirectory() returned file with invalid prefix: %s", file)
					}
				}
			}
		})
	}
}

func TestGetFile(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		wantErr  bool
	}{
		{
			name:     "Empty path",
			filepath: "",
			wantErr:  true,
		},
		{
			name:     "Kind config",
			filepath: "kind.yaml",
			wantErr:  false,
		},
		{
			name:     "Nonexistent file",
			filepath: "nonexistent.txt",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := GetFile(tt.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(content) == 0 {
				t.Error("GetFile() returned empty content for valid file")
			}
		})
	}
}

func TestGetFilesByExtension(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		ext     string
		wantErr bool
	}{
		{
			name:    "YAML files in components",
			dir:     "components",
			ext:     "yaml",
			wantErr: false,
		},
		{
			name:    "YAML files in components with dot",
			dir:     "components",
			ext:     ".yaml",
			wantErr: false,
		},
		{
			name:    "Invalid directory",
			dir:     "nonexistent",
			ext:     "yaml",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := GetFilesByExtension(tt.dir, tt.ext)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFilesByExtension() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify all files have the correct extension
				for _, file := range files {
					if !strings.HasSuffix(file, "."+strings.TrimPrefix(tt.ext, ".")) {
						t.Errorf("GetFilesByExtension() returned file with wrong extension: %s", file)
					}
				}
			}
		})
	}
}

// Helper function to compare slices regardless of order
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aCopy := make([]string, len(a))
	bCopy := make([]string, len(b))
	copy(aCopy, a)
	copy(bCopy, b)
	sort.Strings(aCopy)
	sort.Strings(bCopy)
	return reflect.DeepEqual(aCopy, bCopy)
}
