package helm

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sh31k30ps/gikopsctl/pkg/config/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessHelmChart(t *testing.T) {
	tests := []struct {
		name string
		cfg  *component.HelmChart
	}{
		{
			name: "sans prefix",
			cfg: &component.HelmChart{
				Chart:   "test-chart",
				Version: "1.0.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processHelmChart("test", tt.cfg, "")
		})
	}
}

func TestGetHelmDefaultFile(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *component.HelmChart
		prefix   string
		expected string
	}{
		{
			name: "sans prefix",
			cfg: &component.HelmChart{
				Chart:   "test-chart",
				Version: "1.0.0",
			},
			prefix:   "",
			expected: filepath.Join("default", "values.yaml"),
		},
		{
			name: "avec prefix",
			cfg: &component.HelmChart{
				Chart:   "test-chart",
				Version: "1.0.0",
			},
			prefix:   "dev",
			expected: filepath.Join("default", "dev-values.yaml"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getHelmDefaultFile(tt.cfg, tt.prefix)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProcessHelmDefaults(t *testing.T) {
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

	tests := []struct {
		name        string
		cfg         *component.HelmChart
		prefix      string
		mockCommand func() error // Pour simuler la commande helm
		wantErr     bool
	}{
		{
			name: "Chart inexistant",
			cfg: &component.HelmChart{
				Chart:   "test-chart",
				Version: "1.0.0",
			},
			prefix:  "",
			wantErr: true,
		},
		{
			name: "Chart existant",
			cfg: &component.HelmChart{
				Chart:   "traefik/traefik",
				Version: "v34.4.0",
			},
			prefix:  "",
			wantErr: false,
		},
		{
			name: "Chart existant avec prefix",
			cfg: &component.HelmChart{
				Chart:   "traefik/traefik-crds",
				Version: "v1.4.0",
			},
			prefix:  "crds",
			wantErr: false,
		},
		{
			name: "fichier existe déjà",
			cfg: &component.HelmChart{
				Chart:   "test-chart",
				Version: "1.0.0",
			},
			prefix:  "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockCommand != nil {
				err := tt.mockCommand()
				require.NoError(t, err)
			}

			err := processHelmDefaults(tt.cfg, tt.prefix)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProcessHelmTemplate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *component.HelmChart
		prefix  string
		wantErr bool
	}{
		{
			name: "Chart inexistant",
			cfg: &component.HelmChart{
				Chart:   "test-chart",
				Version: "1.0.0",
			},
			prefix:  "",
			wantErr: true,
		},
		{
			name: "Chart existant",
			cfg: &component.HelmChart{
				Chart:   "traefik/traefik",
				Version: "v34.4.0",
			},
			prefix:  "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := processHelmTemplate(tt.name, tt.cfg, tt.prefix)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProcessCutTemplate(t *testing.T) {
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

	tests := []struct {
		name     string
		content  string
		wantErr  bool
		validate func(t *testing.T)
	}{
		{
			name: "document yaml simple",
			content: `---
# Source: charts/test/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
---`,
			wantErr: false,
			validate: func(t *testing.T) {
				// Vérifier que le fichier a été créé correctement
				content, err := os.ReadFile(filepath.Join("base", "templates", "deployment.yaml"))
				require.NoError(t, err)
				assert.Contains(t, string(content), "kind: Deployment")
			},
		},
		{
			name: "documents multiples",
			content: `---
# Source: charts/test/templates/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: test-svc
---
# Source: charts/test/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
---`,
			wantErr: false,
			validate: func(t *testing.T) {
				// Vérifier les deux fichiers
				svcContent, err := os.ReadFile(filepath.Join("base", "templates", "service.yaml"))
				require.NoError(t, err)
				assert.Contains(t, string(svcContent), "kind: Service")

				deployContent, err := os.ReadFile(filepath.Join("base", "templates", "deployment.yaml"))
				require.NoError(t, err)
				assert.Contains(t, string(deployContent), "kind: Deployment")
			},
		},
		{
			name: "gestion des doublons",
			content: `---
# Source: charts/test/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test1
---
# Source: charts/test/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test2
---`,
			wantErr: false,
			validate: func(t *testing.T) {
				// Vérifier que les deux fichiers existent avec des noms différents
				files, err := filepath.Glob(filepath.Join("base", "templates", "deployment*.yaml"))
				require.NoError(t, err)
				assert.Len(t, files, 2)
			},
		},
		{
			name: "gestion des fichiers vides",
			content: `---
# Source: charts/test/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test1
---
# Source: charts/test/templates/deployment-2.yaml



---
# Source: charts/test/templates/deployment-3.yaml

---
# Source: charts/test/templates/deployment-4.yaml
---`,
			wantErr: false,
			validate: func(t *testing.T) {
				// Vérifier que les deux fichiers existent avec des noms différents
				files, err := filepath.Glob(filepath.Join("base", "templates", "*.yaml"))
				require.NoError(t, err)
				assert.Len(t, files, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Nettoyer le répertoire base avant chaque test
			os.RemoveAll(filepath.Join(tmpDir, "base"))

			// Créer le fichier temporaire avec le contenu de test
			tmpFile := filepath.Join(tmpDir, "test.yaml")
			err := os.WriteFile(tmpFile, []byte(tt.content), 0644)
			require.NoError(t, err)

			// Exécuter le test
			err = processCutTemplate(tmpFile)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Exécuter les validations spécifiques au test
			if tt.validate != nil {
				tt.validate(t)
			}
		})
	}
}

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

			result := getUniqueFilePath(tt.inputPath)
			assert.Equal(t, filepath.Join(tmpDir, tt.expectedSuffix), result)
		})
	}
}

func TestFileExists(t *testing.T) {
	// Créer un répertoire temporaire pour les tests
	tmpDir, err := os.MkdirTemp("", "helm-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Créer un fichier de test
	testFile := filepath.Join(tmpDir, "test.yaml")
	err = os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "fichier existe",
			path:     testFile,
			expected: true,
		},
		{
			name:     "fichier n'existe pas",
			path:     filepath.Join(tmpDir, "nonexistent.yaml"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := false
			if _, err := os.Stat(tt.path); err != nil {
				result = true
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}
