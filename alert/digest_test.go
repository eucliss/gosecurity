package alert

import (
	"os"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create a temporary YAML file for testing
	content := []byte(`
source: test_source
query: test_query
conditions:
  - field: value
    operator: GREATER
    value: "5"
index: test_index
`)
	tmpfile, err := os.CreateTemp("", "test_alert_*.yaml")
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

	// Test Load function
	alert, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	expected := Alert{
		Source: "test_source",
		Query:  "test_query",
		Conditions: []Condition{
			{Field: "value", Operator: "GREATER", Value: "5"},
		},
		Index: "test_index",
	}

	if !reflect.DeepEqual(alert, expected) {
		t.Errorf("Load() = %v, want %v", alert, expected)
	}
}

// func TestCondition_CompareInt(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		condition Condition
// 		first     int
// 		second    int
// 		want      bool
// 	}{
// 		{"Greater true", Condition{Operator: "GREATER"}, 10, 5, true},
// 		{"Greater false", Condition{Operator: "GREATER"}, 5, 10, false},
// 		{"Unknown operator", Condition{Operator: "UNKNOWN"}, 10, 5, false},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := tt.condition.CompareInt(tt.first, tt.second); got != tt.want {
// 				t.Errorf("Condition.CompareInt() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestCondition_CheckInt(t *testing.T) {
	tests := []struct {
		name      string
		condition Condition
		input     map[string]interface{}
		want      bool
	}{
		{
			"Greater true",
			Condition{Field: "value", Operator: "GREATER", Value: "5"},
			map[string]interface{}{"value": "10"},
			true,
		},
		{
			"Less true",
			Condition{Field: "value", Operator: "LESS", Value: "10"},
			map[string]interface{}{"value": "5"},
			true,
		},
		{
			"Equals true",
			Condition{Field: "value", Operator: "EQUALS", Value: "5"},
			map[string]interface{}{"value": "5"},
			true,
		},
		{
			"Invalid field type",
			Condition{Field: "value", Operator: "GREATER", Value: "5"},
			map[string]interface{}{"value": 10}, // int instead of string
			false,
		},
		{
			"Invalid condition value",
			Condition{Field: "value", Operator: "GREATER", Value: "invalid"},
			map[string]interface{}{"value": "10"},
			false,
		},
		{
			"Unknown operator",
			Condition{Field: "value", Operator: "UNKNOWN", Value: "5"},
			map[string]interface{}{"value": "10"},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.condition.CheckInt(tt.input); got != tt.want {
				t.Errorf("Condition.CheckInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlert_Check(t *testing.T) {
	tests := []struct {
		name   string
		alert  Alert
		input  []map[string]interface{}
		want   bool
		checks int
	}{
		{
			"Single condition, all rows meet",
			Alert{
				Conditions: []Condition{{Field: "value", Operator: "GREATER", Value: "5"}},
			},
			[]map[string]interface{}{
				{"value": "10"},
				{"value": "15"},
			},
			true,
			2,
		},
		{
			"Single condition, one row doesn't meet",
			Alert{
				Conditions: []Condition{{Field: "value", Operator: "GREATER", Value: "5"}},
			},
			[]map[string]interface{}{
				{"value": "10"},
				{"value": "3"},
			},
			false,
			1,
		},
		{
			"Multiple conditions, all meet",
			Alert{
				Conditions: []Condition{
					{Field: "value", Operator: "GREATER", Value: "5"},
					{Field: "value", Operator: "EQUALS", Value: "10"},
				},
			},
			[]map[string]interface{}{
				{"value": "10", "status": "active"},
			},
			true,
			1,
		},
		{
			"Multiple conditions, one doesn't meet",
			Alert{
				Conditions: []Condition{
					{Field: "value", Operator: "GREATER", Value: "5"},
					{Field: "value", Operator: "EQUALS", Value: "3"},
				},
			},
			[]map[string]interface{}{
				{"value": "10", "status": "inactive"},
			},
			false,
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.alert.Check(tt.input)
			if got != tt.want {
				t.Errorf("Alert.Check() = %v, want %v", got, tt.want)
			}
			// You might need to add a way to check the number of checks performed
			// This depends on how you implement the counting of checks in your Alert.Check() method
		})
	}
}

func TestCondition_Check(t *testing.T) {
	tests := []struct {
		name      string
		condition Condition
		input     map[string]interface{}
		want      bool
	}{
		{
			"Integer check",
			Condition{Field: "value", Operator: "GREATER", Value: "5"},
			map[string]interface{}{"value": "10"},
			true,
		},
		{
			"Unknown type",
			Condition{Field: "value", Operator: "GREATER", Value: "5"},
			map[string]interface{}{"value": []int{1, 2, 3}},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.condition.Check(tt.input); got != tt.want {
				t.Errorf("Condition.Check() = %v, want %v", got, tt.want)
			}
		})
	}
}
