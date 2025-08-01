package helm

import (
	"fmt"
	"os"
	"testing"

	"github.com/sh31k30ps/gikops/pkg/config/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelm(t *testing.T) {
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

	repo := struct {
		Repo        string
		Url         string
		Version     string
		Chart       string
		CrdsVersion string
		CrdsChart   string
	}{
		Repo:        "traefik",
		Url:         "https://traefik.github.io/charts",
		Version:     "v34.4.0",
		Chart:       "traefik/traefik",
		CrdsVersion: "v1.4.0",
		CrdsChart:   "traefik/traefik-crds",
	}

	tests := []struct {
		cfg     *component.Component
		wantErr bool
		keep    bool
	}{
		{
			cfg: &component.Component{
				Name: "NoHelmConfig",
			},
			wantErr: false,
			keep:    false,
		},
		{
			cfg: &component.Component{
				Name: "MissingRepoName",
				Helm: &component.HelmConfig{
					URL: "http://toto.too.com",
					Chart: &component.HelmChart{
						Chart:   "NoOne",
						Version: "not",
					},
				},
			},
			wantErr: true,
			keep:    false,
		},
		{
			cfg: &component.Component{
				Name: "MissingURL",
				Helm: &component.HelmConfig{
					Repo: "noOne",
					Chart: &component.HelmChart{
						Chart:   "NoOne",
						Version: "not",
					},
				},
			},
			wantErr: true,
			keep:    false,
		},
		{
			cfg: &component.Component{
				Name: "MissingChart",
				Helm: &component.HelmConfig{
					Repo: "noOne",
					URL:  "http://toto.too.com",
				},
			},
			wantErr: true,
			keep:    false,
		},
		{
			cfg: &component.Component{
				Name: "RepoFailed",
				Helm: &component.HelmConfig{
					Repo: "failed",
					URL:  "https://toto.too.com",
					Chart: &component.HelmChart{
						Chart:   "NoOne",
						Version: "not",
					},
				},
			},
			wantErr: true,
			keep:    false,
		},
		{
			cfg: &component.Component{
				Name: "ChartFailed",
				Helm: &component.HelmConfig{
					Repo: repo.Repo,
					URL:  repo.Url,
					Chart: &component.HelmChart{
						Chart:   "NoOne",
						Version: "not",
					},
				},
			},
			wantErr: true,
			keep:    false,
		},
		{
			cfg: &component.Component{
				Name: "chart-existant",
				Helm: &component.HelmConfig{
					Repo: repo.Repo,
					URL:  repo.Url,
					Chart: &component.HelmChart{
						Chart:   repo.Chart,
						Version: repo.Version,
					},
				},
			},
			wantErr: false,
			keep:    false,
		},
		{
			cfg: &component.Component{
				Name: "crds-other-chart",
				Helm: &component.HelmConfig{
					Repo: repo.Repo,
					URL:  repo.Url,
					Chart: &component.HelmChart{
						Chart:   repo.Chart,
						Version: repo.Version,
					},
					CRDsChart: &component.HelmChart{
						Chart:   repo.CrdsChart,
						Version: repo.CrdsVersion,
					},
				},
			},
			wantErr: false,
			keep:    false,
		},
		{
			cfg: &component.Component{
				Name: "keep-tmp",
				Helm: &component.HelmConfig{
					Repo: repo.Repo,
					URL:  repo.Url,
					Chart: &component.HelmChart{
						Chart:   repo.Chart,
						Version: repo.Version,
					},
				},
			},
			wantErr: false,
			keep:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.cfg.Name, func(t *testing.T) {
			err := setupHelmRepo(tt.cfg.Name, tt.cfg, tt.keep)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.cfg.Helm != nil {
					tmpFile := fmt.Sprintf("%s/default/values.yaml", tmpDir)
					assert.FileExists(t, tmpFile, "defaults values must be generated")
					assert.FileExists(t, fmt.Sprintf("%s/base/kustomization.yaml", tmpDir))
					if tt.cfg.Helm.CRDsChart != nil {
						assert.FileExists(t, fmt.Sprintf("%s/default/crds-values.yaml", tmpDir))
					}
					if tt.keep {
						assert.FileExists(t, fmt.Sprintf("%s/base/tmp.yaml", tmpDir))
					} else {
						assert.NoFileExists(t, fmt.Sprintf("%s/base/tmp.yaml", tmpDir))
					}
				}
			}
		})
	}
}
