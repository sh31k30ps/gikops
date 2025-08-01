package helm

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateKustomizeFile_WhenFileDoesNotExist(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "helm-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Save the original working directory and switch to tmp
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Create the "base" folder and some dummy files inside
	baseDir := filepath.Join(tmpDir, "base")
	err = os.MkdirAll(baseDir, 0755)
	require.NoError(t, err)

	dummyFile := filepath.Join(baseDir, "dummy.yaml")
	err = os.WriteFile(dummyFile, []byte("content"), 0644)
	require.NoError(t, err)

	// Now test the function
	err = createKustomizeFile()
	require.NoError(t, err)
}

func TestCreateKustomizeFile_WhenFileAlreadyExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "helm-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	baseDir := filepath.Join(tmpDir, "base")
	err = os.MkdirAll(baseDir, 0755)
	require.NoError(t, err)

	kustomFile := filepath.Join(baseDir, "kustomization.yaml")
	err = os.WriteFile(kustomFile, []byte("existing"), 0644)
	require.NoError(t, err)

	// Should do nothing and return no error
	err = createKustomizeFile()
	require.NoError(t, err)
}
