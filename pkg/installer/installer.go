package installer

import (
	"fmt"
	"os/exec"

	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
	"github.com/sh31k30ps/gikopsctl/pkg/internal/cli"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cmd"
)

type Installer struct {
	config   *project.Project
	provider *cluster.Provider
	logger   log.Logger
	status   *cli.Status
}

func NewInstaller(config *project.Project, logger log.Logger) *Installer {
	var provider cluster.ProviderOption
	switch config.LocalCluster.KindConfig.Provider {
	case "docker":
		provider = cluster.ProviderWithDocker()
	case "podman":
		provider = cluster.ProviderWithPodman()
	case "nerdctl":
		provider = cluster.ProviderWithNerdctl("nerdctl")
	default:
		provider = cluster.ProviderWithDocker()
	}

	return &Installer{
		config: config,
		provider: cluster.NewProvider(
			cluster.ProviderWithLogger(cmd.NewLogger()),
			provider,
		),
		logger: logger,
		status: cli.StatusForLogger(logger),
	}
}

func (i *Installer) Install() error {
	i.logger.V(0).Info("Starting installation")

	// Create kind network
	i.status.Start("Creating kind network")
	if err := i.createKindNetwork(); err != nil {
		i.status.End(false)
		return fmt.Errorf("failed to create kind network: %w", err)
	}
	i.status.End(true)

	// Create kind cluster
	if err := i.createCluster(); err != nil {
		return fmt.Errorf("failed to create cluster: %w", err)
	}

	i.logger.V(0).Info("Overrides")
	for _, override := range i.config.LocalCluster.KindConfig.OverridesFolder {
		i.status.Start(fmt.Sprintf("Applying override: %s", override))
		if err := i.applyOverridesFolder(override); err != nil {
			i.status.End(false)
			return fmt.Errorf("failed to apply override: %w", err)
		}
		i.status.End(true)
	}
	return nil
}

func (i *Installer) Uninstall() error {
	i.logger.V(0).Info("Starting uninstallation")

	// Delete kind cluster
	if err := i.provider.Delete(i.clusterName(), ""); err != nil {
		return fmt.Errorf("failed to delete cluster: %w", err)
	}
	i.status.Start("Deleting kind cluster")
	i.status.End(true)

	// Remove kind network
	i.status.Start("Removing kind network")
	if err := i.removeKindNetwork(); err != nil {
		i.status.End(false)
		return fmt.Errorf("failed to remove kind network: %w", err)
	}
	i.status.End(true)

	return nil
}

func (i *Installer) createKindNetwork() error {
	cmd := exec.Command("docker", "network", "create",
		"--driver=bridge",
		"--subnet=192.168.100.0/24",
		"--gateway=192.168.100.1",
		"kind")
	// cmd.Stdout = nil
	// cmd.Stderr = nil
	return cmd.Run()
}

func (i *Installer) removeKindNetwork() error {
	cmd := exec.Command("docker", "network", "rm", "kind")
	// cmd.Stdout = nil
	// cmd.Stderr = nil
	return cmd.Run()
}

func (i *Installer) createCluster() error {
	options := []cluster.CreateOption{
		cluster.CreateWithConfigFile(i.config.LocalCluster.KindConfig.ConfigFile),
	}

	return i.provider.Create(i.clusterName(), options...)
}

func (i *Installer) applyOverridesFolder(folder string) error {
	cmd := exec.Command("kubectl", "apply", "-f", folder)
	if err := i.logger.V(1).CmdOutput(cmd); err != nil {
		return err
	}

	return nil
}

// func (i *Installer) applyCoreDNS() error {
// 	cmd := exec.Command("kubectl", "apply", "-f", "./overrides/coreDns.yaml")
// 	if err := i.logger.V(1).CmdOutput(cmd); err != nil {
// 		return err
// 	}

// 	cmd = exec.Command("kubectl", "rollout", "restart", "deployment/coredns", "-n", "kube-system")
// 	if err := i.logger.V(1).CmdOutput(cmd); err != nil {
// 		return err
// 	}

// 	cmd = exec.Command("kubectl", "rollout", "status", "deployment/coredns", "-n", "kube-system")
// 	return i.logger.V(1).CmdOutput(cmd)
// }

func (i *Installer) clusterName() string {
	return i.config.GetLocalClusterName()
}
