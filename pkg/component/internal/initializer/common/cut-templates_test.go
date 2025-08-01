package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			err = CutTemplate(tmpFile)
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
