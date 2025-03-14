package component

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sh31k30ps/gikopsctl/pkg/config/cluster"
)

const (
	errMesgExists                = "cluster already exists"
	errMesgLocalFolder           = "local folder is not a directory"
	errMesgClusterFolderNotFound = "cluster folder not found"
)

func (m *Manager) AddCluster(cmpt string, cl cluster.Cluster) error {
	clFolder := filepath.Join(cmpt, cl.Name())
	clLocal := filepath.Join(cmpt, "local")
	if _, err := os.Stat(clFolder); err == nil {
		return errors.New(errMesgExists)
	}
	if err := os.MkdirAll(clFolder, 0755); err != nil {
		return err
	}

	if flocal, err := os.Stat(clLocal); err == nil {
		if !flocal.IsDir() {
			return errors.New(errMesgLocalFolder)
		}
		err := filepath.Walk(clLocal, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			relPath, err := filepath.Rel(clLocal, path)
			if err != nil {
				return err
			}
			destPath := filepath.Join(clFolder, relPath)
			if info.IsDir() {
				return os.MkdirAll(destPath, info.Mode())
			}
			srcFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer srcFile.Close()
			destFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer destFile.Close()
			_, err = io.Copy(destFile, srcFile)
			return err
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) DeleteCluster(cpmt string, cl cluster.Cluster) error {
	clFolder := filepath.Join(cpmt, cl.Name())
	if _, err := os.Stat(clFolder); err != nil {
		return errors.New(errMesgClusterFolderNotFound)
	}
	return os.RemoveAll(clFolder)
}

func IsErrorClusterFolderExists(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), errMesgExists)
}

func IsErrorLocalFolder(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), errMesgLocalFolder)
}

func IsErrorClusterFolderNotFound(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), errMesgClusterFolderNotFound)
}
