package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

var ConfigPath string
var SourcesPath string

func LoadSources() (SourcesConfig, error) {
	fmt.Println("Loading sources config")
	if SourcesPath == "" {
		SourcesPath = "sources.yaml"
	}
	file, err := os.Open(SourcesPath)
	if err != nil {
		return SourcesConfig{}, err
	}
	defer file.Close()

	var config SourcesConfig
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return SourcesConfig{}, err
	}
	return config, err
}

func LoadConfig() (MonitoringConfig, error) {
	fmt.Println("Loading base config")
	if ConfigPath == "" {
		ConfigPath = "config.yaml"
	}
	file, err := os.Open(ConfigPath)
	if err != nil {
		return MonitoringConfig{}, fmt.Errorf("error reading YAML file: %w", err)
	}
	defer file.Close()

	var config MonitoringConfig

	// Unmarshal the YAML data into the Config struct
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return MonitoringConfig{}, err
	}
	return config, err
}
