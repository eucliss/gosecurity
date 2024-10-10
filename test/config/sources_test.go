package securitygo

import (
	"fmt"
	"gosecurity/config"
	"gosecurity/indexer"
	"os"
	"testing"
)

var DummySimple config.FileSource
var DummyJson config.FileSource

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
	// Run all the tests
	exitCode := m.Run()

	// Exit with the code returned by m.Run()
	os.Exit(exitCode)
}

func TestLoadFileSimple(t *testing.T) {
	fmt.Println("Testing LoadFile")
	data, err := indexer.LoadFile(DummySimple)
	if err != nil {
		t.Errorf("Error loading file: %v", err)
	}
	if len(data) == 0 {
		t.Errorf("Expected file data, got empty data")
	}
	if string(data) != "Hello World!" {
		t.Errorf("Expected 'test data', got %v", string(data))
	}

}

func TestLoadFileJSON(t *testing.T) {
	fmt.Println("Testing LoadFile")
	data, err := indexer.LoadFile(DummyJson)
	if err != nil {
		t.Errorf("Error loading file: %v", err)
	}
	j, err := indexer.ReadJSONFile(data)
	if len(j) == 0 {
		t.Errorf("Expected file data, got empty data")
	}
	if err != nil {
		t.Errorf("Error reading JSON: %v", err)
	}
	// j[0]
	if j[0]["destination_ip"] != "10.0.0.5" {
		t.Errorf("Expcted J[0] to have destination_ip 10.0.0.5")
	}
	if j[1]["source_ip"] != "192.168.1.150" {
		t.Errorf("Expcted J[1] to have source_ip of 192.168.1.150")
	}
	if j[2]["status"] != "ALLOW" {
		t.Errorf("Expected J[2] to have status of ALLOW")
	}
}
