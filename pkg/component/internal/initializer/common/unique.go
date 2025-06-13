package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// getUniqueFilePath retourne un chemin de fichier unique en ajoutant un suffixe numérique si nécessaire
func UniqueFilePath(basePath string) string {
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return basePath
	}

	ext := filepath.Ext(basePath)
	baseWithoutExt := strings.TrimSuffix(basePath, ext)

	counter := 1
	for {
		newPath := fmt.Sprintf("%s_%d%s", baseWithoutExt, counter, ext)
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
		counter++
	}
}
