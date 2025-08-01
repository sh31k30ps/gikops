package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUniqueFilePath(t *testing.T) {
	// Créer un répertoire temporaire pour les tests
	tmpDir, err := os.MkdirTemp("", "helm-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name           string
		setupFiles     []string
		inputPath      string
		expectedSuffix string
	}{
		{
			name:           "fichier n'existe pas",
			setupFiles:     []string{},
			inputPath:      filepath.Join(tmpDir, "test.yaml"),
			expectedSuffix: "test.yaml",
		},
		{
			name:           "fichier existe une fois",
			setupFiles:     []string{"test.yaml"},
			inputPath:      filepath.Join(tmpDir, "test.yaml"),
			expectedSuffix: "test_1.yaml",
		},
		{
			name:           "fichier existe plusieurs fois",
			setupFiles:     []string{"test.yaml", "test_1.yaml"},
			inputPath:      filepath.Join(tmpDir, "test.yaml"),
			expectedSuffix: "test_2.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Créer les fichiers de test
			for _, f := range tt.setupFiles {
				err := os.WriteFile(filepath.Join(tmpDir, f), []byte("test"), 0644)
				require.NoError(t, err)
			}

			result := UniqueFilePath(tt.inputPath)
			assert.Equal(t, filepath.Join(tmpDir, tt.expectedSuffix), result)
		})
	}
}
