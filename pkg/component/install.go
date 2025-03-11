package component

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sh31k30ps/gikopsctl/assets"
	"github.com/sh31k30ps/gikopsctl/pkg/config/project"
)

func InstallComponents(components project.ProjectComponent) error {
	if err := os.MkdirAll(components.Name, 0755); err != nil {
		return err
	}
	if components.Require != nil {
		for _, require := range components.Require {
			if err := Install(components.Name, require); err != nil {
				return err
			}
		}
	}
	return nil
}

func Install(base, component string) error {
	files, err := assets.GetComponentFiles(component)
	if err != nil {
		return err
	}
	for _, file := range files {
		content, err := assets.GetFile(file)
		if err != nil {
			return err
		}
		parts := strings.Split(file, "/")
		dest := path.Join(base, path.Join(parts[2:]...))
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(dest, content, 0644); err != nil {
			return err
		}
	}
	return nil
}
