package config

import (
	"strings"
)

type FileSource struct {
	Name        string `yaml:"name"`
	Path        string `yaml:"path"`
	Format      string `yaml:"format"`
	Description string `yaml:"description"`
	Source      string `yaml:"source"`
	TargetIndex string `yaml:"target_index"`
}

type SourcesConfig struct {
	// Config represents the overall configuration structure.
	FileSources []FileSource `yaml:"file_sources"`
}

func (s SourcesConfig) toString() (finalString string) {
	// Example array of strings

	// Build a string from the array
	var result strings.Builder // Using strings.Builder for efficiency
	result.WriteString("File Sources:\n")
	for _, file := range s.FileSources {
		result.WriteString(file.Name + "\n")
		result.WriteString(file.Path + "\n")
		result.WriteString(file.Description + "\n")
		result.WriteString("---\n")
	}

	// Convert the builder to a string
	finalString = result.String()
	return
}

func (s SourcesConfig) String() string {
	return s.toString()
}
