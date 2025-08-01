package basiccreator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sh31k30ps/gikopsctl/pkg/config"
	"github.com/sh31k30ps/gikopsctl/pkg/config/cluster"
	"github.com/sh31k30ps/gikopsctl/pkg/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreator_Create(t *testing.T) {
	// Create a mock logger
	logger := services.NewLogger("test")

	// Create a new Creator instance
	creator := NewCreator(logger)

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "basiccreator-test-*")
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
		// Create a BasicCluster config
		clusterName := "test-cluster"
		cfg := cluster.NewBasicCluster()
		cfg.SetName(clusterName)

		// Call the Create function
		err = creator.Create(cfg)

		// Assert that the error is nil
		require.NoError(t, err, "Create should not return an error")

		// Assert that the folder was created
		folderPath := filepath.Join("clusters", clusterName)
		_, err = os.Stat(folderPath)
		assert.NoError(t, err, "Folder should be created")

	})

	// Test case 2: Error when config is not BasicCluster
	t.Run("Error when config is not BasicCluster", func(t *testing.T) {
		var cfg config.ConfigObject
		// Call the Create function
		err = creator.Create(cfg)

		// Assert that the error is not nil
		assert.Error(t, err, "Create should return an error")
		assert.EqualError(t, err, "config is not a BasicCluster", "Error message should match")
	})
}
