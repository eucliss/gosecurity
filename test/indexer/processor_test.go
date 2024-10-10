package gosecurity

import (
	"gosecurity/indexer"
	"testing"
)

func TestParseCSV(t *testing.T) {
	res, err := indexer.ProcessCSV(DummyCSV.Path)
	if err != nil {
		t.Errorf("Error loading file: %v", err)
	}
	if res[0]["machine_id"] != "Tractor_1" {
		t.Errorf("Incorrect machine_id for machine1")
	}

	if res[2]["fuel_consumption_liters"] != "2" {
		t.Errorf("Incorrect fuel consumption for machine 3")
	}
}

func TestLoadFileJSON(t *testing.T) {
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
