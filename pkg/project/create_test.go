package project

import (
	"os"
	"testing"

	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
	"github.com/sh31k30ps/gikopsctl/pkg/services"
)

func TestProjectCreator_Create(t *testing.T) {
	// Créer un répertoire temporaire pour les tests
	tmpDir, err := os.MkdirTemp("", "project-creator-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Changer le répertoire de travail
	originalWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalWd)

	// Créer une instance de ProjectCreator
	logger := services.NewLogger("test")
	pc := NewProjectCreator(logger)

	// Configurer les données de test
	testConfig := project.NewProject()
	testConfig.Name = "test"
	cluster := project.NewKindCluster()
	cluster.SetName("test")
	testConfig.Clusters = []project.ProjectCluster{
		cluster,
	}
	testConfig.Components = []project.ProjectComponent{
		{
			Name: "core",
			Require: []string{
				"ingress/traefik",
				"monitoring/prometheus",
				"security/cert-manager",
				"security/mkcert",
			},
		},
	}
	// Exécuter le test
	err = pc.Create(testConfig)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	// Vérifier que les fichiers attendus ont été créés
	expectedFiles := []string{
		".gitignore",
		"clusters/test/kind.yaml",
		"clusters/test/overrides/coreDns.yaml",
		"core/traefik/gikcpnt.yaml",
		"core/prometheus/gikcpnt.yaml",
		"core/cert-manager/gikcpnt.yaml",
		"core/mkcert/gikcpnt.yaml",
	}

	for _, file := range expectedFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not created", file)
		}
	}
}
