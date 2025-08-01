package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sh31k30ps/gikopsctl/pkg/config/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestBaseDir(t *testing.T, files []string) string {
	baseDir := filepath.Join("base")
	err := os.MkdirAll(baseDir, 0755)
	assert.NoError(t, err)

	for _, file := range files {
		path := filepath.Join(baseDir, file)
		err := os.WriteFile(path, []byte("test content"), 0644)
		assert.NoError(t, err)
	}

	return baseDir
}

func TestCleanBaseDir_NoPreserve(t *testing.T) {
	// Créer un répertoire temporaire pour les tests
	tmpDir, err := os.MkdirTemp("", "helm-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Sauvegarder le répertoire courant
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	// Changer vers le répertoire temporaire
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	setupTestBaseDir(t, []string{"file1.txt", "file2.txt", "kustomization.yaml"})

	cfg := &component.Component{}

	err = CleanBaseDir(cfg)
	assert.NoError(t, err)

	files, _ := os.ReadDir("base")
	assert.Equal(t, 1, len(files))
	assert.Equal(t, "kustomization.yaml", files[0].Name())
}

func TestCleanBaseDir_WithPreserve(t *testing.T) {
	// Créer un répertoire temporaire pour les tests
	tmpDir, err := os.MkdirTemp("", "helm-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Sauvegarder le répertoire courant
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	// Changer vers le répertoire temporaire
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	setupTestBaseDir(t, []string{"file1.txt", "keepme.yaml", "kustomization.yaml"})

	cfg := &component.Component{
		Files: &component.FilesConfig{
			Keep: []string{"keepme.yaml"},
		},
	}

	err = CleanBaseDir(cfg)
	assert.NoError(t, err)

	files, _ := os.ReadDir("base")
	names := []string{}
	for _, f := range files {
		names = append(names, f.Name())
	}

	assert.ElementsMatch(t, []string{"kustomization.yaml", "keepme.yaml"}, names)
}

func TestCleanBaseDir_CreateIfMissing(t *testing.T) {
	// Créer un répertoire temporaire pour les tests
	tmpDir, err := os.MkdirTemp("", "helm-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Sauvegarder le répertoire courant
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	// Changer vers le répertoire temporaire
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	cfg := &component.Component{}
	err = CleanBaseDir(cfg)
	assert.NoError(t, err)

	_, statErr := os.Stat("base")
	assert.NoError(t, statErr)
}
