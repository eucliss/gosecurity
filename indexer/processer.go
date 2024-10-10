package indexer

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gosecurity/config"
	"os"
)

func LoadFile(f config.FileSource) ([]byte, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		return data, fmt.Errorf("error reading YAML file: %w", err)
	}
	return data, nil
}

func ReadJSONFile(b []byte) ([]map[string]interface{}, error) {
	// Unmarshal the JSON data into a slice of maps
	var res []map[string]interface{}
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// Eventually maybe allow for configs to be rendered to parse certain sources
func ProcessJSON(v []byte, targetIndex string) {
	j, _ := ReadJSONFile(v)
	fmt.Println("JSON")
	fmt.Println(j)
	for index := range j {
		fmt.Println("Processing index: ", index)
		go Send(targetIndex, j[index])
	}
}

func Send(targetIndex string, e map[string]interface{}) {
	NormalizedChannel <- NormalizedEvent{Event: e, Index: targetIndex}
}

func ProcessCSV(filePath string, targetIndex string) ([]map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var output []map[string]interface{}
	headers := records[0]
	for _, record := range records[1:] {
		row := make(map[string]interface{})
		for i, value := range record {
			row[headers[i]] = value
		}
		go Send(targetIndex, row)
		output = append(output, row)
	}
	return output, nil
}

func ProcessXML(v []byte, targetIndex string) {

	var result []map[string]interface{}
	err := xml.Unmarshal(v, &result)
	fmt.Println(err)
	fmt.Println(result)

	// fmt.Println("new approach")
	// file, _ := os.Open("indexer/data/crop.xml")
	// byteValue, _ := ioutil.ReadAll(file)
	// fmt.Println(byteValue)
	// Convert result into []map[string]interface{}
	for i := range result {
		go Send(targetIndex, result[i])
	}

}

// func parseXML(filePath string) ([]map[string]interface{}, error) {
// 	// Read XML file
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	byteValue, _ := os.ReadAll(file)

// 	var result map[string]interface{}
// 	err = xml.Unmarshal(byteValue, &result)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Convert result into []map[string]interface{}
// 	output := []map[string]interface{}{result}
// 	return output, nil
// }
