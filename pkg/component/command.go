package component

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sh31k30ps/gikopsctl/assets"
	"github.com/sh31k30ps/gikopsctl/pkg/cli"
	"github.com/sh31k30ps/gikopsctl/pkg/config/cluster"
	"github.com/sh31k30ps/gikopsctl/pkg/config/component"
	"github.com/sh31k30ps/gikopsctl/pkg/config/manager"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
	"github.com/sh31k30ps/gikopsctl/pkg/directories"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/sh31k30ps/gikopsctl/pkg/services"
	uicomponent "github.com/sh31k30ps/gikopsctl/pkg/ui/component"
	"github.com/sh31k30ps/gikopsctl/pkg/ui/standard"
)

type Command struct {
	manager *Manager
	logger  log.Logger
	status  *cli.Status
	ui      *uicomponent.UIComponentRequester
}

func NewCommand(logger log.Logger) *Command {
	return &Command{
		logger:  logger,
		status:  cli.StatusForLogger(logger),
		ui:      uicomponent.NewRequester(logger),
		manager: NewManager(logger),
	}
}

func (c *Command) Create(args ...interface{}) error {
	if len(args) < 1 {
		return fmt.Errorf("component folder is required")
	}
	folder := args[0].(string)

	project, err := services.GetCurrentProject()
	if err != nil {
		return err
	}

	if project.GetComponent(folder) == nil {
		return fmt.Errorf("Component folder %s not found", folder)
	}

	if _, err := c.ui.Request(""); err != nil {
		return err
	}

	cfg, err := c.ui.Config()
	if err != nil {
		return err
	}
	cfgC := cfg.(*component.Component)
	if err := manager.SaveComponent(filepath.Join(folder, cfgC.Name), cfgC); err != nil {
		return err
	}

	for _, cl := range project.Clusters {
		if err := c.AddCluster(filepath.Join(folder, cfgC.Name), cl); err != nil {
			return err
		}
	}
	c.logger.V(0).Info("Component created successfully")
	return nil
}

func (c *Command) Edit(mode string, args ...interface{}) error {
	return nil
}

func (c *Command) Delete(id interface{}) error {
	return nil
}

func (c *Command) Add(mode string, args ...string) error {
	return nil
}

func (c *Command) Install() error {
	cfg, err := services.GetCurrentProject()
	if err != nil {
		return err
	}
	cpmtRoots := cfg.Components
	for _, cpmtRoot := range cpmtRoots {
		if err := c.InstallRoot(cpmtRoot); err != nil {
			return err
		}
	}
	return nil
}

func (c *Command) InstallRoot(cpmtRoot project.ProjectComponent) error {
	if err := os.MkdirAll(cpmtRoot.Name, 0755); err != nil {
		return err
	}
	if cpmtRoot.Require != nil {
		for _, require := range cpmtRoot.Require {
			if err := c.InstallComponent(cpmtRoot.Name, require); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Command) InstallComponent(cpmtRootName, cpmtName string) error {
	files, err := assets.GetComponentFiles(cpmtName)
	if err != nil {
		return err
	}
	sp := strings.Split(cpmtName, "/")
	cpmtSimpleName := sp[len(sp)-1]
	cpmtSimpleFolder := filepath.Join(cpmtRootName, cpmtSimpleName)

	if _, err := os.Stat(cpmtSimpleFolder); err == nil {
		overwrite, err := standard.PromptYesNo(fmt.Sprintf("Folder %s already exists, overwrite?", cpmtSimpleFolder))
		if err != nil {
			return err
		}
		if !overwrite {
			return nil
		}
		if err := os.RemoveAll(cpmtSimpleFolder); err != nil {
			return err
		}
	}

	for _, file := range files {
		content, err := assets.GetFile(file)
		if err != nil {
			return err
		}
		parts := strings.Split(file, "/")
		dest := filepath.Join(cpmtRootName, filepath.Join(parts[2:]...))

		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(dest, content, 0644); err != nil {
			return err
		}
	}
	return nil
}

func (c *Command) AddCluster(cmpt string, cl cluster.Cluster) error {
	if err := c.manager.AddCluster(cmpt, cl); err != nil {
		if !IsErrorClusterFolderExists(err) && !IsErrorLocalFolder(err) {
			return err
		}
		return nil
	}
	c.logger.V(1).Info(fmt.Sprintf("Cluster %s added to %s", cl.Name(), cmpt))
	return nil
}

func (c *Command) CleanComponentsCluster(cl cluster.Cluster) error {
	project, err := services.GetCurrentProject()
	if err != nil {
		return err
	}
	components := directories.GetRootsComponents(project)
	mngr := NewManager(c.logger)
	for _, cmpt := range components {
		if err := mngr.DeleteCluster(cmpt, cl); err != nil {
			if !IsErrorClusterFolderNotFound(err) {
				return err
			}
			continue
		}
		c.logger.V(1).Info(fmt.Sprintf("Cluster %s removed from %s", cl.Name(), cmpt))
	}
	return nil
}
