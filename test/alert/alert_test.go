package gosecurity

import (
	"gosecurity/alert"
	"os"
	"testing"
)

var DummyPath string

func TestMain(m *testing.M) {
	DummyPath = "./dummy.yaml"
	// Run all the tests
	exitCode := m.Run()

	// Exit with the code returned by m.Run()
	os.Exit(exitCode)
}

func TestLoadAlertConfig(t *testing.T) {
	alert, err := alert.Load(DummyPath)
	if err != nil {
		t.Errorf("Error loading alert config: %v", err)
	}
	if alert.Source != "dummy" {
		t.Errorf("Expected 'dummy', got %v", alert.Source)
	}
	if alert.Query == "" {
		t.Errorf("Expected query, got %v", alert.Query)
	}
}
