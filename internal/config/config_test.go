package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// Save original environment variables
	originalNode := os.Getenv("NODE")
	originalPrivateKeys := os.Getenv("PRIVATE_KEYS")

	// Clean up after test
	defer func() {
		os.Setenv("NODE", originalNode)
		os.Setenv("PRIVATE_KEYS", originalPrivateKeys)
		viper.Reset()
	}()

	// Set environment variables for testing
	os.Setenv("NODE", "http://localhost:8545")
	os.Setenv("PRIVATE_KEYS", "key1,key2")

	config, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "http://localhost:8545", config.Node)
	assert.Equal(t, []string{"key1", "key2"}, config.PrivateKeys)
	assert.Equal(t, "./main.db", config.DatabasePath) // default value
}

func TestLoadWithCustomDatabasePath(t *testing.T) {
	// Save original environment variables
	originalNode := os.Getenv("NODE")
	originalPrivateKeys := os.Getenv("PRIVATE_KEYS")
	originalDatabasePath := os.Getenv("DATABASE_PATH")

	// Clean up after test
	defer func() {
		os.Setenv("NODE", originalNode)
		os.Setenv("PRIVATE_KEYS", originalPrivateKeys)
		os.Setenv("DATABASE_PATH", originalDatabasePath)
		viper.Reset()
	}()

	// Set environment variables for testing
	os.Setenv("NODE", "http://localhost:8545")
	os.Setenv("PRIVATE_KEYS", "key1,key2")
	os.Setenv("DATABASE_PATH", "./custom.db")

	config, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "http://localhost:8545", config.Node)
	assert.Equal(t, []string{"key1", "key2"}, config.PrivateKeys)
	assert.Equal(t, "./custom.db", config.DatabasePath)
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Node:         "http://localhost:8545",
				PrivateKeys:  []string{"key1"},
				DatabasePath: "./test.db",
			},
			wantErr: false,
		},
		{
			name: "missing node",
			config: Config{
				Node:         "",
				PrivateKeys:  []string{"key1"},
				DatabasePath: "./test.db",
			},
			wantErr: true,
		},
		{
			name: "missing private keys",
			config: Config{
				Node:         "http://localhost:8545",
				PrivateKeys:  []string{},
				DatabasePath: "./test.db",
			},
			wantErr: true,
		},
		{
			name: "missing database path",
			config: Config{
				Node:         "http://localhost:8545",
				PrivateKeys:  []string{"key1"},
				DatabasePath: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoadMissingConfigFiles(t *testing.T) {
	// Save original environment variables
	originalNode := os.Getenv("NODE")
	originalPrivateKeys := os.Getenv("PRIVATE_KEYS")

	// Clean up after test
	defer func() {
		os.Setenv("NODE", originalNode)
		os.Setenv("PRIVATE_KEYS", originalPrivateKeys)
		viper.Reset()
	}()

	// Set environment variables for testing
	os.Setenv("NODE", "http://localhost:8545")
	os.Setenv("PRIVATE_KEYS", "key1,key2")

	// Test loading config when no config files exist
	config, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "http://localhost:8545", config.Node)
	assert.Equal(t, []string{"key1", "key2"}, config.PrivateKeys)
	assert.Equal(t, "./main.db", config.DatabasePath) // default value
}
