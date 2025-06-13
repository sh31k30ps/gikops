package services

import (
	"strings"
)

func GetComponentFolderFromDepth() string {
	cfg, err := GetCurrentProject()
	if err != nil {
		return ""
	}
	if cfg.Level == 0 {
		return ""
	}
	dirs := strings.Split(cfg.Origin, "/")
	dir := dirs[len(dirs)-cfg.Level]
	if cfg.GetComponent(dir) == nil {
		return ""
	}
	return dir
}

func GetComponentFromDepth() string {
	cfg, err := GetCurrentProject()
	if err != nil {
		return ""
	}

	if cfg.Level < 2 {
		return ""
	}

	if dir := GetComponentFolderFromDepth(); dir == "" {
		return ""
	}

	dirs := strings.Split(cfg.Origin, "/")
	index := len(dirs) - (3 - cfg.Level)
	if index < 0 {
		return ""
	}
	return dirs[index]
}
