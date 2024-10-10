package gosecurity

import (
	"gosecurity/config"
	"os"
	"testing"
)

var DummySimple config.FileSource
var DummyJson config.FileSource
var DummyCSV config.FileSource

func TestMain(m *testing.M) {
	DummySimple = config.FileSource{
		Name:        "Simple",
		Format:      "txt",
		Path:        "../data/simple.txt",
		Description: "Simple text",
	}
	DummyJson = config.FileSource{
		Name:        "Network",
		Format:      "json",
		Path:        "../data/network.json",
		Description: "JSON file",
	}
	DummyCSV = config.FileSource{
		Name:        "Machinery",
		Format:      "csv",
		Path:        "../data/machinery.csv",
		Description: "CSV File",
	}
	// Run all the tests
	exitCode := m.Run()

	// Exit with the code returned by m.Run()
	os.Exit(exitCode)
}
