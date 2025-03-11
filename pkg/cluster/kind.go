package cluster

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sh31k30ps/gikopsctl/pkg/cli"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"

	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"sigs.k8s.io/kind/pkg/cluster"
	kcmd "sigs.k8s.io/kind/pkg/cmd"
	klog "sigs.k8s.io/kind/pkg/log"
)

type KindInstaller struct {
	logger log.Logger
	status *cli.Status
}

func NewKindInstaller(logger log.Logger) *KindInstaller {
	return &KindInstaller{
		logger: logger,
		status: cli.StatusForLogger(logger),
	}
}

func (i *KindInstaller) Install(c project.ProjectCluster) error {
	cfg := c.Config().(*project.KindConfig)
	provider := getKindProvider(cfg)
	i.logger.V(0).Info("Starting installation")

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	os.Chdir(filepath.Join("clusters", c.Name()))
	defer os.Chdir(currentDir)

	// Create kind network
	i.status.Start("Creating kind network")
	if err := createKindNetwork(); err != nil {
		i.status.End(false)
		return fmt.Errorf("failed to create kind network: %w", err)
	}
	i.status.End(true)

	// Create kind cluster
	options := []cluster.CreateOption{
		cluster.CreateWithConfigFile(c.Config().(*project.KindConfig).ConfigFile),
	}

	if err := provider.Create(c.GetClusterName(), options...); err != nil {
		return fmt.Errorf("failed to create cluster: %w", err)
	}

	i.logger.V(0).Info("Overrides")
	for _, override := range cfg.OverridesFolder {
		i.status.Start(fmt.Sprintf("Applying override: %s", override))
		if err := i.applyOverridesFolder(override); err != nil {
			i.status.End(false)
			return fmt.Errorf("failed to apply override: %w", err)
		}
		i.status.End(true)
	}
	return nil
}

func (i *KindInstaller) Uninstall(c project.ProjectCluster) error {
	i.logger.V(0).Info("Starting uninstallation")
	cfg := c.Config().(*project.KindConfig)
	// Delete kind cluster
	if err := getKindProvider(cfg).Delete(c.GetClusterName(), ""); err != nil {
		return fmt.Errorf("failed to delete cluster: %w", err)
	}
	i.status.Start("Deleting kind cluster")
	i.status.End(true)

	// Remove kind network
	i.status.Start("Removing kind network")
	if err := removeKindNetwork(); err != nil {
		i.status.End(false)
		return fmt.Errorf("failed to remove kind network: %w", err)
	}
	i.status.End(true)

	return nil
}

func (i *KindInstaller) applyOverridesFolder(folder string) error {
	cmd := exec.Command("kubectl", "apply", "-f", folder)
	if err := i.logger.V(1).CmdOutput(cmd); err != nil {
		return err
	}

	return nil
}

func createKindNetwork() error {
	cmd := exec.Command("docker", "network", "create",
		"--driver=bridge",
		"--subnet=192.168.100.0/24",
		"--gateway=192.168.100.1",
		"kind")
	// cmd.Stdout = nil
	// cmd.Stderr = nil
	return cmd.Run()
}

func removeKindNetwork() error {
	cmd := exec.Command("docker", "network", "rm", "kind")
	// cmd.Stdout = nil
	// cmd.Stderr = nil
	return cmd.Run()
}

func getKindProvider(cfg *project.KindConfig) *cluster.Provider {
	var provider cluster.ProviderOption
	switch cfg.Provider {
	case "docker":
		provider = cluster.ProviderWithDocker()
	case "podman":
		provider = cluster.ProviderWithPodman()
	case "nerdctl":
		provider = cluster.ProviderWithNerdctl("nerdctl")
	default:
		provider = cluster.ProviderWithDocker()
	}
	return cluster.NewProvider(provider, cluster.ProviderWithLogger(
		getKindLogger(),
	))
}

func getKindLogger() klog.Logger {
	return kcmd.NewLogger()
}
