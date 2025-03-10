package component

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/sh31k30ps/gikopsctl/pkg/config/component"
)

func setupHelmRepo(name string, cfg *component.Component) error {
	if cfg.Helm == nil || cfg.Helm.Repo == "" || cfg.Helm.URL == "" {
		return nil
	}

	if err := cleanBaseDir(cfg); err != nil {
		return fmt.Errorf("failed to clean base directory: %w", err)
	}

	addCmd := exec.Command("helm", "repo", "add", cfg.Helm.Repo, cfg.Helm.URL)
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("failed to add helm repo: %w", err)
	}

	updateCmd := exec.Command("helm", "repo", "update")
	if err := updateCmd.Run(); err != nil {
		return fmt.Errorf("failed to update helm repos: %w", err)
	}
	cmdArgs := []string{}
	if cfg.Helm.CRDsChart == nil || cfg.Helm.CRDsChart.Chart == "" {
		cmdArgs = append(cmdArgs, "--include-crds")
	}

	namespace := component.GetComponentPrefix(name)
	if cfg.Namespace != "" {
		namespace = cfg.Namespace
	}
	cmdArgs = append(cmdArgs, "-n", namespace)

	if err := processHelmChart("", cfg.Helm.CRDsChart, "crds"); err != nil {
		return err
	}

	if err := processHelmChart(cfg.Name, cfg.Helm.Chart, "", cmdArgs...); err != nil {
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

func processHelmChart(name string, chart *component.HelmChart, prefix string, args ...string) error {
	if chart == nil || chart.Chart == "" {
		return nil
	}
	if err := processHelmDefaults(chart, prefix); err != nil {
		return err
	}
	if err := processHelmTemplate(name, chart, prefix, args...); err != nil {
		return err
	}
	return nil
}

func getHelmDefaultFile(cfg *component.HelmChart, prefix string) string {
	defaultsFile := "values.yaml"
	if prefix != "" {
		defaultsFile = prefix + "-" + defaultsFile
	}

	return filepath.Join("default", defaultsFile)
}

func processHelmDefaults(cfg *component.HelmChart, prefix string) error {
	if cfg.Chart == "" {
		return nil
	}
	defaultsPath := getHelmDefaultFile(cfg, prefix)
	if _, err := os.Stat(defaultsPath); err == nil {
		return nil
	}

	cmdArgs := []string{"show", "values", cfg.Chart}
	if cfg.Version != "" {
		cmdArgs = append(cmdArgs, "--version", cfg.Version)
	}
	showCmd := exec.Command("helm", cmdArgs...)
	output, err := showCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get helm values: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(defaultsPath), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(defaultsPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write values file: %w", err)
	}

	return nil
}

func processHelmTemplate(name string, cfg *component.HelmChart, prefix string, options ...string) error {
	if cfg.Version == "" || cfg.Chart == "" {
		return nil
	}
	cmdArgs := []string{"template"}
	if name != "" {
		cmdArgs = append(cmdArgs, name)
	}
	cmdArgs = append(cmdArgs, "--values", getHelmDefaultFile(cfg, prefix))

	if len(options) > 0 {
		cmdArgs = append(cmdArgs, options...)
	}
	if cfg.Version != "" {
		cmdArgs = append(cmdArgs, "--version", cfg.Version)
	}
	cmdArgs = append(cmdArgs, cfg.Chart)

	if err := os.MkdirAll("base", 0755); err != nil {
		return err
	}

	templateCmd := exec.Command("helm", cmdArgs...)
	output, err := templateCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get helm templates: %w", err)
	}
	if prefix != "" {
		file := filepath.Join("base", prefix+".yaml")
		if err := os.WriteFile(file, output, 0644); err != nil {
			return fmt.Errorf("failed to write values file: %w", err)
		}
		return nil
	}

	tmpFile := filepath.Join("base", "tmp.yaml")
	if err := os.WriteFile(tmpFile, output, 0644); err != nil {
		return fmt.Errorf("failed to write values file: %w", err)
	}

	if err := processCutTemplate(tmpFile); err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	os.Remove(tmpFile)

	return nil
}

func processCutTemplate(file string) error {
	// Lire le contenu du fichier
	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Diviser le contenu en documents YAML (séparés par ---)
	documents := strings.Split(string(content), "---\n")

	for _, doc := range documents {
		if strings.TrimSpace(doc) == "" {
			continue
		}

		lines := strings.Split(doc, "\n")
		var sourceLine string
		var sourceIndex int

		// Chercher la ligne "# Source:"
		for i, line := range lines {
			if strings.HasPrefix(line, "# Source: ") {
				sourceLine = line
				sourceIndex = i
				break
			}
		}

		if sourceLine == "" {
			continue
		}

		// Extraire et traiter le chemin
		originalPath := strings.TrimPrefix(sourceLine, "# Source: ")
		pathParts := strings.Split(originalPath, "/")

		// Construire le nouveau chemin en ignorant les 2 premiers niveaux
		var pathBuilder strings.Builder
		for i := 2; i < len(pathParts); i++ {
			if pathBuilder.Len() > 0 {
				pathBuilder.WriteString("/")
			}
			pathBuilder.WriteString(pathParts[i])
		}

		path := pathBuilder.String()
		if path == "" {
			log.Printf("Warning: Path is empty after removing levels: %s", originalPath)
			continue
		}

		// Préparer le contenu du nouveau fichier
		var contentBuilder strings.Builder
		contentBuilder.WriteString("---\n")

		lineCount := 1
		// Ajouter les lignes pertinentes
		for i := sourceIndex + 1; i < len(lines); i++ {
			line := lines[i]
			if !strings.HasPrefix(line, "---") && !strings.HasPrefix(line, "# Source: ") && line != "" {
				lineCount++
				contentBuilder.WriteString(line)
				contentBuilder.WriteString("\n")
			}
		}

		if lineCount == 1 {
			continue
		}

		// Construire le chemin complet du fichier
		filePath := filepath.Join("base", path)
		dirPath := filepath.Dir(filePath)

		// Créer la structure de répertoires
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			log.Printf("Error creating directory structure: %s", dirPath)
			continue
		}

		// Gérer les fichiers en doublon
		finalPath := getUniqueFilePath(filePath)

		// Écrire le fichier
		if err := os.WriteFile(finalPath, []byte(contentBuilder.String()), 0644); err != nil {
			log.Printf("Error writing file %s: %v", finalPath, err)
			continue
		}
	}

	return nil
}

// getUniqueFilePath retourne un chemin de fichier unique en ajoutant un suffixe numérique si nécessaire
func getUniqueFilePath(basePath string) string {
	if !fileExists(basePath) {
		return basePath
	}

	ext := filepath.Ext(basePath)
	baseWithoutExt := strings.TrimSuffix(basePath, ext)

	counter := 1
	for {
		newPath := fmt.Sprintf("%s_%d%s", baseWithoutExt, counter, ext)
		if !fileExists(newPath) {
			return newPath
		}
		counter++
	}
}

func cleanBaseDir(cfg *component.Component) error {
	// Clean up base directory except kustomization.yaml and configured files
	baseDir := filepath.Join("base")
	files, err := os.ReadDir(baseDir)
	if err != nil {
		return fmt.Errorf("failed to read base directory: %w", err)
	}

	// Keep track of files that should be preserved
	preserveFiles := []string{
		"kustomization.yaml",
	}

	if cfg.Files != nil && len(cfg.Files.Keep) > 0 {
		preserveFiles = append(preserveFiles, cfg.Files.Keep...)
	}

	// Remove files that aren't in the preserve list
	for _, f := range files {
		if !slices.Contains(preserveFiles, f.Name()) {
			filePath := filepath.Join(baseDir, f.Name())
			if f.IsDir() {
				if err := os.RemoveAll(filePath); err != nil {
					return fmt.Errorf("failed to remove directory %s: %w", f.Name(), err)
				}
			} else {
				if err := os.Remove(filePath); err != nil {
					return fmt.Errorf("failed to remove file %s: %w", f.Name(), err)
				}
			}
		}
	}
	return nil
}

// fileExists vérifie si un fichier existe
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
