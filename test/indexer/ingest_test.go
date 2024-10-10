package gosecurity

import (
	"gosecurity/indexer"
	"testing"
)

func TestLoadFileSimple(t *testing.T) {
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
