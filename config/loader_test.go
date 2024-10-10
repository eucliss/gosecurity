package config

import (
	"os"
	"reflect"
	"testing"
)

var osOpen = os.Open

func TestLoadSources(t *testing.T) {
	// Create a temporary YAML file for testing
	content := []byte(`
file_sources:
  - name: source1
    path: http://example.com/source1
  - name: source2
    path: http://example.com/source2
`)
	tmpfile, err := os.CreateTemp("", "sources.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Temporarily replace the original file with our test file

	// ... existing code ...
	originalOpen := os.Open
	osOpen = func(name string) (*os.File, error) {
		return os.Open(tmpfile.Name())
	}
	defer func() { osOpen = originalOpen }()

	// Temporarily replace the config file path
	oldConfigPath := SourcesPath // Assume ConfigPath is exported from the loader package
	SourcesPath = tmpfile.Name()
	defer func() { SourcesPath = oldConfigPath }()

	// Test the LoadSources function
	config, err := LoadSources()
	if err != nil {
		t.Fatalf("LoadSources() error = %v", err)
	}

	expected := SourcesConfig{
		FileSources: []FileSource{
			{Name: "source1", Path: "http://example.com/source1"},
			{Name: "source2", Path: "http://example.com/source2"},
		},
	}

	if !reflect.DeepEqual(config, expected) {
		t.Errorf("LoadSources() = %v, want %v", config, expected)
	}
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file for testing
	tempFile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test configuration to the temp file
	testConfig := []byte(`
server:
  host: "localhost"
  port: "8080"
database:
  type: "elastic"
  cert_location: "db/http_ca.crt"
  host: "localhost"
  port: "9200"
  user: "elastic"
  password: "testpassword"
  new_index_on_launch: "test-index"
`)
	if _, err := tempFile.Write(testConfig); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Temporarily replace the config file path
	oldConfigPath := ConfigPath // Assume ConfigPath is exported from the loader package
	ConfigPath = tempFile.Name()
	defer func() { ConfigPath = oldConfigPath }()

	// Test the LoadConfig function
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Add assertions to check if the config was loaded correctly
	if config.Server.Host != "localhost" {
		t.Errorf("Expected server host to be 'localhost', got %s", config.Server.Host)
	}
	// Add more assertions as needed
}

// TestLoadSourcesFileNotFound tests the error case when the sources file is not found
func TestLoadSourcesFileNotFound(t *testing.T) {
	// Ensure the file doesn't exist
	os.Remove("sources.yaml")

	_, err := LoadSources()
	if err == nil {
		t.Error("LoadSources() expected an error, got nil")
	}
}

// TestLoadConfigFileNotFound tests the error case when the config file is not found
func TestLoadConfigFileNotFound(t *testing.T) {
	// Ensure the file doesn't exist
	os.Remove("config.yaml")

	_, err := LoadConfig()
	if err == nil {
		t.Error("LoadConfig() expected an error, got nil")
	}
}
