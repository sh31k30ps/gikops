package helm

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/sh31k30ps/gikops/pkg/component/internal/initializer/common"
	"github.com/sh31k30ps/gikops/pkg/config/component"
)

func setupHelmRepo(name string, cfg *component.Component, keepTmp bool) error {
	if cfg.Helm == nil {
		return nil
	}

	if cfg.Helm.Repo == "" || cfg.Helm.URL == "" {
		return fmt.Errorf("missing basic helm configurations: %v", []string{"repo", "url"})
	}

	if err := common.CleanBaseDir(cfg); err != nil {
		return fmt.Errorf("failed to clean base directory: %w", err)
	}

	if strings.HasPrefix(cfg.Helm.URL, "https://") {
		addCmd := exec.Command("helm", "repo", "add", cfg.Helm.Repo, cfg.Helm.URL)
		if err := addCmd.Run(); err != nil {
			return fmt.Errorf("failed to add helm repo: %w", err)
		}

		updateCmd := exec.Command("helm", "repo", "update")
		if err := updateCmd.Run(); err != nil {
			return fmt.Errorf("failed to update helm repos: %w", err)
		}
	}

	cmdArgs := []string{}
	if cfg.Helm.CRDsChart == nil || cfg.Helm.CRDsChart.Chart == "" {
		cmdArgs = append(cmdArgs, "--include-crds")
	} else {
		if err := processHelmChart(cfg.Name, cfg.Helm.CRDsChart, "crds", false); err != nil {
			return err
		}
	}

	namespace := component.GetComponentPrefix(name)
	if cfg.Namespace != "" {
		namespace = cfg.Namespace
	}
	cmdArgs = append(cmdArgs, "-n", namespace)

	if err := processHelmChart(cfg.Name, cfg.Helm.Chart, "", keepTmp, cmdArgs...); err != nil {
		return err
	}

	return nil
}

func processHookInit(cfg *component.HelmInitHooks) error {
	if cfg == nil {
		return nil
	}

	for _, upload := range cfg.Uploads {
		if err := handleUpload(upload); err != nil {
			return err
		}
	}

	for _, resolve := range cfg.Resolves {
		if err := handleResolve(resolve); err != nil {
			return err
		}
	}

	for _, rename := range cfg.Renames {
		if err := handleRename(rename); err != nil {
			return err
		}
	}

	for _, contact := range cfg.Concats {
		if err := handleConcat(contact); err != nil {
			return err
		}
	}

	return nil
}

func processHelmChart(name string, chart *component.HelmChart, prefix string, keepTmp bool, args ...string) error {
	if chart == nil || chart.Chart == "" {
		return fmt.Errorf("missing chart name config")
	}
	if err := processHelmDefaults(chart, prefix); err != nil {
		return err
	}
	if err := processHelmTemplate(name, chart, prefix, keepTmp, args...); err != nil {
		return err
	}
	if prefix != "" {
		return nil
	}
	if err := createKustomizeFile(); err != nil {
		return err
	}
	return nil
}
