package kustomize

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sh31k30ps/gikopsctl/pkg/cli"
	"github.com/sh31k30ps/gikopsctl/pkg/component/internal/initializer/common"
	"github.com/sh31k30ps/gikopsctl/pkg/config/component"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/sh31k30ps/gikopsctl/pkg/services"
)

type Initializer struct {
	logger log.Logger
	status *cli.Status
}

func NewInitializer(logger log.Logger, status *cli.Status) *Initializer {
	return &Initializer{
		logger: logger,
		status: status,
	}
}

func (i *Initializer) Init(name string, keepTmp bool) error {
	cCfg, err := services.GetComponent(name)
	if err != nil {
		return fmt.Errorf("failed to get component: %w", err)
	}
	if cCfg.Kustomize != nil {
		return setupKuztomizeRepo(name, cCfg, keepTmp)
	}
	return nil
}

func setupKuztomizeRepo(name string, cfg *component.Component, keepTmp bool) error {
	if cfg.Kustomize == nil || len(cfg.Kustomize.URLs) == 0 {
		return nil
	}
	baseDir := filepath.Join("base")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return err
	}
	for _, url := range cfg.Kustomize.URLs {
		dest := filepath.Join(baseDir, filepath.Base(url))
		if err := getURLContent(url, dest); err != nil {
			return fmt.Errorf("failed to download Kustomize repo: %w", err)
		}
		if err := common.CutTemplate(dest); err != nil {
			return fmt.Errorf("failed to cut templates: %w", err)
		}
	}
	if err := createKustomizeFile(); err != nil {
		return err
	}
	return nil
}

func getURLContent(url string, dest string) error {
	content, err := fetchURL(url)
	if err != nil {
		return err
	}
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.Write(content); err != nil {
		return err
	}
	return nil
}

func fetchURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
