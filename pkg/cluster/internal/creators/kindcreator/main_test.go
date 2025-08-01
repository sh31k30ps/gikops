package kindcreator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sh31k30ps/gikops/pkg/config"
	"github.com/sh31k30ps/gikops/pkg/config/cluster"
	"github.com/sh31k30ps/gikops/pkg/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreator_Create(t *testing.T) {
	// Create a mock logger
	logger := services.NewLogger("test")

	// Create a new Creator instance
	creator := NewCreator(logger)

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "kindcreator-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Change to the temporary directory
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Test case 1: Successful creation
	t.Run("Successful creation", func(t *testing.T) {
		// Create a KindCluster config
		clusterName := "test-kind-cluster"
		cfg := cluster.NewKindCluster()
		cfg.SetName(clusterName)

		// Call the Create function
		err = creator.Create(cfg)

		// Assert that the error is nil
		require.NoError(t, err, "Create should not return an error")

		// Assert that the folder was created
		folderPath := filepath.Join("clusters", clusterName)
		_, err = os.Stat(folderPath)
		assert.NoError(t, err, "Folder should be created")

		// Assert that kind.yaml exists
		kindConfigPath := filepath.Join(folderPath, "kind.yaml")
		_, err = os.Stat(kindConfigPath)
		assert.NoError(t, err, "kind.yaml should be created")
	})

	// Test case 2: Error when config is not KindCluster
	t.Run("Error when config is not KindCluster", func(t *testing.T) {
		var cfg config.ConfigObject
		// Call the Create function
		err = creator.Create(cfg)

		// Assert that the error is not nil
		assert.Error(t, err, "Create should return an error")
		assert.EqualError(t, err, "config is not a KindCluster", "Error message should match")
	})

	// Test case 3: Test copyKindOverrides
	t.Run("Test copyKindOverrides", func(t *testing.T) {
		clusterName := "test-kind-cluster-overrides"
		cfg := cluster.NewKindCluster()
		cfg.SetName(clusterName)
		kindConfig := cfg.Config().(*cluster.KindConfig)
		kindConfig.OverridesFolder = []string{"overrides"} // Assuming coreDns is a valid folder in assets

		folderPath := filepath.Join("clusters", clusterName)
		err := os.MkdirAll(folderPath, 0755)
		require.NoError(t, err)

		err = creator.copyKindOverrides(folderPath, cfg)
		require.NoError(t, err, "copyKindOverrides should not return an error")

		// Assert that the coreDns folder exists
		overridesFolderPath := filepath.Join(folderPath, "overrides")
		_, err = os.Stat(overridesFolderPath)
		assert.NoError(t, err, "overrides folder should be created")

		// Assert that a file within coreDns exists (assuming there is a file in assets/overrides/coreDns.yaml)
		coreDnsFilePath := filepath.Join(overridesFolderPath, "coreDns.yaml")
		_, err = os.Stat(coreDnsFilePath)
		assert.NoError(t, err, "coreDns.yaml should be created")

		os.RemoveAll(folderPath)
	})

	// // Test case 4: Error creating folder
	// t.Run("Error creating folder", func(t *testing.T) {
	// 	clusterName := "test-kind-cluster-error"
	// 	cfg := cluster.NewKindCluster()
	// 	cfg.SetName(clusterName)
	// 	// folderPath := filepath.Join("clusters", clusterName)

	// 	// Simulate an error by making the parent directory unwritable
	// 	err := os.MkdirAll("clusters", 0444) // Read-only for everyone
	// 	require.NoError(t, err, "Failed to create read-only clusters dir")

	// 	// Call the Create function, which should fail to create the cluster folder
	// 	err = creator.Create(cfg)

	// 	// Assert that an error is returned
	// 	assert.Error(t, err, "Create should return an error")

	// 	// Clean up: remove the read-only clusters directory
	// 	os.RemoveAll("clusters")
	// })
}
